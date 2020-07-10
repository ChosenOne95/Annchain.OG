package types

import (
	"fmt"
	"github.com/annchain/OG/arefactor/og/types"
	"github.com/annchain/OG/common"
	"github.com/annchain/OG/common/byteutil"
	"github.com/annchain/OG/common/hexutil"
	"strings"
)

// Sequencer is the block generated by consensus to confirm the graph
type Sequencer struct {
	// graph structure info
	Hash         types.Hash
	ParentsHash  types.Hashes
	Height       uint64
	MineNonce    uint64
	AccountNonce uint64
	Issuer       common.Address
	Signature    hexutil.Bytes
	PublicKey    hexutil.Bytes
	StateRoot    types.Hash
	//Proposing    bool `msg:"-"` // is the sequencer is proposal ,did't commit yet ,use this flag to avoid bls sig verification failed

	// derived properties
	Weight  uint64
	invalid bool
}

func (s *Sequencer) SetMineNonce(v uint64) {
	s.MineNonce = v
}

func (s *Sequencer) SetParents(hashes types.Hashes) {
	s.ParentsHash = hashes
}

func (s *Sequencer) SetWeight(weight uint64) {
	s.Weight = weight
}

func (s *Sequencer) SetValid(b bool) {
	s.invalid = !b
}

func (s *Sequencer) Valid() bool {
	return !s.invalid
}

func (s *Sequencer) SetSender(addr common.Address) {
	s.Issuer = addr
}

func (s *Sequencer) SetHash(h types.Hash) {
	s.Hash = h
}

func (s *Sequencer) GetNonce() uint64 {
	return s.AccountNonce
}

func (s *Sequencer) Sender() common.Address {
	return s.Issuer
}

func (s *Sequencer) SetHeight(height uint64) {
	s.Height = height
}

func (s *Sequencer) Dump() string {
	var phashes []string
	for _, p := range s.ParentsHash {
		phashes = append(phashes, p.Hex())
	}
	return fmt.Sprintf("pHash:[%s], Issuer : %s , Height: %d, nonce : %d , blspub: %s, signatute : %s, pubkey:  %s root: %s",
		strings.Join(phashes, " ,"),
		s.Issuer.Hex(),
		s.Height,
		s.AccountNonce,
		s.PublicKey,
		hexutil.Encode(s.PublicKey),
		hexutil.Encode(s.Signature),
		s.StateRoot.Hex(),
	)
}



func (s *Sequencer) SignatureTargets() []byte {
	w := byteutil.NewBinaryWriter()

	w.Write(s.PublicKey, s.AccountNonce)
	w.Write(s.Issuer.Bytes)

	//w.Write(s.Height, s.Weight, s.StateRoot.KeyBytes)
	w.Write(s.Height, s.StateRoot.Bytes)
	for _, parent := range s.GetParents() {
		w.Write(parent.Bytes)
	}
	return w.Bytes()
}

func (s *Sequencer) GetType() TxBaseType {
	return TxBaseTypeSequencer
}

func (s *Sequencer) GetHeight() uint64 {
	return s.Height
}

func (s *Sequencer) GetWeight() uint64 {
	if s.Weight == 0 {
		panic("implementation error: weight not initialized")
	}
	return s.Weight
}

func (s *Sequencer) GetHash() types.Hash {
	return s.Hash
}

func (s *Sequencer) GetParents() types.Hashes {
	return s.ParentsHash
}

func (s *Sequencer) String() string {
	return fmt.Sprintf("Sq-[%.10s]:%d", s.Issuer.String(), s.AccountNonce)
	//if s.Issuer == nil {
	//	return fmt.Sprintf("Sq-[nil]-%d", s.AccountNonce)
	//} else {
	//
	//}
}

func (s *Sequencer) CalculateWeight(parents Txis) uint64 {
	var maxWeight uint64
	for _, p := range parents {
		if p.GetWeight() > maxWeight {
			maxWeight = p.GetWeight()
		}
	}
	return maxWeight + 1
}

func (s *Sequencer) Compare(tx Txi) bool {
	switch tx := tx.(type) {
	case *Sequencer:
		if s.GetHash().Cmp(tx.GetHash()) == 0 {
			return true
		}
		return false
	default:
		return false
	}
}