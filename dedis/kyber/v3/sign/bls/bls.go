// Package bls implements the Boneh-Lynn-Shacham (BLS) signature scheme which
// was introduced in the paper "Short Signatures from the Weil Pairing". BLS
// requires pairing-based cryptography.
package bls

import (
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/annchain/OG/dedis/kyber/v3"
	"github.com/annchain/OG/dedis/kyber/v3/pairing"
)

// NewKeyPair creates a new BLS signing key pair. The private key x is a scalar
// and the public key X is a point on curve G2.
func NewKeyPair(suite pairing.Suite, random cipher.Stream) (kyber.Scalar, kyber.Point) {
	x := suite.G2().Scalar().Pick(random)
	X := suite.G2().Point().Mul(x, nil)
	return x, X
}

// Sign creates a BLS signature S = x * H(m) on a message m using the private
// key x. The signature S is a point on curve G1.
func Sign(suite pairing.Suite, x kyber.Scalar, msg []byte) ([]byte, error) {
	HM := hashToPoint(suite, msg)
	xHM := HM.Mul(x, HM)
	s, err := xHM.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// AggregateSignatures combines signatures created using the Sign function
func AggregateSignatures(suite pairing.Suite, sigs ...[]byte) ([]byte, error) {
	sig := suite.G1().Point()
	for _, sigBytes := range sigs {
		sigToAdd := suite.G1().Point()
		if err := sigToAdd.UnmarshalBinary(sigBytes); err != nil {
			return nil, err
		}
		sig.Add(sig, sigToAdd)
	}
	return sig.MarshalBinary()
}

// AggregatePublicKeys takes a slice of public G2 points and returns
// the sum of those points. This is used to verify multisignatures.
func AggregatePublicKeys(suite pairing.Suite, Xs ...kyber.Point) kyber.Point {
	aggregated := suite.G2().Point()
	for _, X := range Xs {
		aggregated.Add(aggregated, X)
	}
	return aggregated
}

// BatchVerify verifies a large number of publicKey/msg pairings with a single aggregated signature.
// Since aggregation is generally much faster than verification, this can be a speed enhancement.
// Benchmarks show a roughly 50% performance increase over individual signature verification
// Every msg must be unique or there is the possibility to accept an invalid signature
// see: https://crypto.stackexchange.com/questions/56288/is-bls-signature-scheme-strongly-unforgeable/56290
// for a description of why each message must be unique.
func BatchVerify(suite pairing.Suite, publics []kyber.Point, msgs [][]byte, sig []byte) error {
	if !distinct(msgs) {
		return fmt.Errorf("bls: error, messages must be distinct")
	}

	s := suite.G1().Point()
	if err := s.UnmarshalBinary(sig); err != nil {
		return err
	}

	var aggregatedLeft kyber.Point
	for i := range msgs {
		hm := hashToPoint(suite, msgs[i])
		pair := suite.Pair(hm, publics[i])

		if i == 0 {
			aggregatedLeft = pair
		} else {
			aggregatedLeft.Add(aggregatedLeft, pair)
		}
	}

	right := suite.Pair(s, suite.G2().Point().Base())
	if !aggregatedLeft.Equal(right) {
		return errors.New("bls: invalid signature")
	}
	return nil
}

// Verify checks the given BLS signature S on the message m using the public
// key X by verifying that the equality e(H(m), X) == e(H(m), x*B2) ==
// e(x*H(m), B2) == e(S, B2) holds where e is the pairing operation and B2 is
// the base point from curve G2.
func Verify(suite pairing.Suite, X kyber.Point, msg, sig []byte) error {
	HM := hashToPoint(suite, msg)
	left := suite.Pair(HM, X)
	s := suite.G1().Point()
	if err := s.UnmarshalBinary(sig); err != nil {
		return err
	}
	right := suite.Pair(s, suite.G2().Point().Base())
	if !left.Equal(right) {
		return errors.New("bls: invalid signature")
	}
	return nil
}

// hashToPoint hashes a message to a point on curve G1. XXX: This should be replaced
// eventually by a proper hash-to-point mapping like Elligator.
func hashToPoint(suite pairing.Suite, msg []byte) kyber.Point {
	h := suite.Hash()
	h.Write(msg)
	x := suite.G1().Scalar().SetBytes(h.Sum(nil))
	return suite.G1().Point().Mul(x, nil)
}

func distinct(msgs [][]byte) bool {
	m := make(map[[32]byte]bool)
	for _, msg := range msgs {
		h := sha256.Sum256(msg)
		if m[h] {
			return false
		}
		m[h] = true
	}
	return true
}
