// Copyright © 2019 Annchain Authors <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package rpc

//go:generate msgp
import (
	"fmt"
	"net/http"
	"time"

	"github.com/annchain/OG/common"
	"github.com/annchain/OG/common/crypto"
	"github.com/annchain/OG/common/math"
	"github.com/annchain/OG/status"
	"github.com/annchain/OG/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	BaseTokenID int32 = 0
)

func (r *RpcController) NewTransaction(c *gin.Context) {
	var (
		tx    types.Txi
		txReq NewTxRequest
		//sig   crypto.Signature
		//pub   crypto.PublicKey
	)

	if status.ArchiveMode {
		Response(c, http.StatusBadRequest, fmt.Errorf("archive mode"), nil)
		return
	}

	err := c.ShouldBindJSON(&txReq)
	if err != nil {
		Response(c, http.StatusBadRequest, fmt.Errorf("request format error: %v", err), nil)
		return
	}

	var parentsHash []common.Hash
	for _, hashStr := range txReq.Parents {
		hash := common.HexToHash(hashStr)
		if hash.Empty() {
			Response(c, http.StatusBadRequest, fmt.Errorf("invalid parent hash: %b", hash.Bytes), nil)
			return
		}
		parentsHash = append(parentsHash, hash)
	}

	from, err := common.StringToAddress(txReq.From)
	if err != nil {
		Response(c, http.StatusBadRequest, fmt.Errorf("from address format error: %v", err), nil)
		return
	}

	// empty address
	to := common.BytesToAddress(nil)
	// 0 bigint
	value := math.NewBigInt(0)

	guarantee, ok := math.NewBigIntFromString(txReq.Guarantee, 10)
	if !ok {
		err = fmt.Errorf("convert guarantee to big int error")
	}
	if err != nil {
		Response(c, http.StatusBadRequest, err, nil)
		return
	}
	//if guarantee.Value.Cmp(big.NewInt(100)) < 0 {
	//	Response(c, http.StatusBadRequest, fmt.Errorf("guarantee should be larger than 100"), nil)
	//	return
	//}

	nonce := txReq.Nonce

	//signature := common.FromHex(txReq.Signature)
	//if signature == nil || txReq.Signature == "" {
	//	Response(c, http.StatusBadRequest, fmt.Errorf("signature format error"), nil)
	//	return
	//}

	//pub, err = crypto.Secp256k1PublicKeyFromString(txReq.Pubkey)
	//if err != nil {
	//	Response(c, http.StatusBadRequest, fmt.Errorf("pubkey format error %v", err), nil)
	//	return
	//}

	//sig = crypto.SignatureFromBytes(pub.Type, signature)
	//if sig.Type != crypto.Signer.GetCryptoType() || pub.Type != crypto.Signer.GetCryptoType() {
	//	Response(c, http.StatusOK, fmt.Errorf("crypto algorithm mismatch"), nil)
	//	return
	//}

	hash := common.HexToHash(txReq.Hash)
	tx, err = r.TxCreator.NewTxForHackathonViewer(from, to, value, guarantee, []byte{}, nonce, parentsHash, BaseTokenID, hash)

	//tx, err = r.TxCreator.NewTxForHackathon(from, to, value, guarantee, []byte{}, nonce, parentsHash, pub, sig, BaseTokenID)
	if err != nil {
		Response(c, http.StatusInternalServerError, fmt.Errorf("new tx failed: %v", err), nil)
		return
	}

	//ok = r.FormatVerifier.VerifySignature(tx)
	//if !ok {
	//	logrus.WithField("request ", txReq).WithField("tx ", tx).Warn("signature invalid")
	//	Response(c, http.StatusInternalServerError, fmt.Errorf("signature invalid"), nil)
	//	return
	//}
	//ok = r.FormatVerifier.VerifySourceAddress(tx)
	//if !ok {
	//	logrus.WithField("request ", txReq).WithField("tx ", tx).Warn("source address invalid")
	//	Response(c, http.StatusInternalServerError, fmt.Errorf("source address invalid"), nil)
	//	return
	//}

	tx.SetVerified(types.VerifiedFormat)
	logrus.WithField("tx", tx).Debugf("tx generated")
	if !r.SyncerManager.IncrementalSyncer.Enabled {
		Response(c, http.StatusOK, fmt.Errorf("tx is disabled when syncing"), nil)
		return
	}

	r.TxBuffer.ReceivedNewTxChan <- tx

	Response(c, http.StatusOK, nil, tx.GetTxHash().Hex())
	return
}

func (r *RpcController) NewSequencer(c *gin.Context) {
	var (
		seq    types.Txi
		seqReq NewSeqRequest
		//sig    crypto.Signature
		//pub    crypto.PublicKey
	)

	err := c.ShouldBindJSON(&seqReq)
	if err != nil {
		Response(c, http.StatusBadRequest, fmt.Errorf("request format error: %v", err), nil)
		return
	}

	var parentsHash []common.Hash
	for _, hashStr := range seqReq.Parents {
		hash := common.HexToHash(hashStr)
		if hash.Empty() {
			Response(c, http.StatusBadRequest, fmt.Errorf("invalid parent hash: %b", hash.Bytes), nil)
			return
		}
		parentsHash = append(parentsHash, hash)
	}

	from, err := common.StringToAddress(seqReq.From)
	if err != nil {
		Response(c, http.StatusBadRequest, fmt.Errorf("from address format error: %v", err), nil)
		return
	}

	treasure, ok := math.NewBigIntFromString(seqReq.Treasure, 10)
	if !ok {
		err = fmt.Errorf("convert guarantee to big int error")
	}
	if err != nil {
		Response(c, http.StatusBadRequest, err, nil)
		return
	}
	//if guarantee.Value.Cmp(big.NewInt(100)) < 0 {
	//	Response(c, http.StatusBadRequest, fmt.Errorf("guarantee should be larger than 100"), nil)
	//	return
	//}

	nonce := seqReq.Nonce

	//signature := common.FromHex(seqReq.Signature)
	//if signature == nil || seqReq.Signature == "" {
	//	Response(c, http.StatusBadRequest, fmt.Errorf("signature format error"), nil)
	//	return
	//}

	//pub, err = crypto.Secp256k1PublicKeyFromString(seqReq.Pubkey)
	//if err != nil {
	//	Response(c, http.StatusBadRequest, fmt.Errorf("pubkey format error %v", err), nil)
	//	return
	//}

	//sig = crypto.SignatureFromBytes(pub.Type, signature)
	//if sig.Type != crypto.Signer.GetCryptoType() || pub.Type != crypto.Signer.GetCryptoType() {
	//	Response(c, http.StatusOK, fmt.Errorf("crypto algorithm mismatch"), nil)
	//	return
	//}

	hash := common.HexToHash(seqReq.Hash)
	height := seqReq.Height
	seq, err = r.TxCreator.NewSeqForHackathonViewer(from, treasure, nonce, parentsHash, hash, height)
	if err != nil {
		Response(c, http.StatusInternalServerError, fmt.Errorf("new seq failed: %v", err), nil)
		return
	}

	seq.SetVerified(types.VerifiedFormat)
	logrus.WithField("seq", seq).Debugf("seq generated")
	if !r.SyncerManager.IncrementalSyncer.Enabled {
		Response(c, http.StatusOK, fmt.Errorf("seq is disabled when syncing"), nil)
		return
	}

	r.TxBuffer.ReceivedNewTxChan <- seq

	Response(c, http.StatusOK, nil, seq.GetTxHash().Hex())
	return
}

//NewTxrequest for RPC request
//msgp:tuple NewTxRequest
type NewTxRequest struct {
	Hash    string   `json:"hash"`
	Parents []string `json:"parents"`
	Nonce   uint64   `json:"nonce"`
	From    string   `json:"from"`
	//To        string   `json:"to"`
	//Value     string   `json:"value"`
	Guarantee string `json:"guarantee"`
	//Data      string `json:"data"`
	//CryptoType string `json:"crypto_type"`
	Signature string `json:"signature"`
	Pubkey    string `json:"pubkey"`
	//TokenId   int32  `json:"token_id"`
}

type NewSeqRequest struct {
	Hash      string   `json:"hash"`
	Parents   []string `json:"parents"`
	Nonce     uint64   `json:"nonce"`
	From      string   `json:"from"`
	Treasure  string   `json:"treasure"`
	Height    uint64   `json:"height"`
	Signature string   `json:"signature"`
	Pubkey    string   `json:"pubkey"`
}

//msgp:tuple NewTxsRequests
type NewTxsRequests struct {
	Txs []NewTxRequest `json:"txs"`
}

func (r *RpcController) NewTransactions(c *gin.Context) {
	var (
		//txs types.Txis
		txrequsets NewTxsRequests
		sig        crypto.Signature
		pub        crypto.PublicKey
		hashes     common.Hashes
	)

	if status.ArchiveMode {
		Response(c, http.StatusBadRequest, fmt.Errorf("archive mode"), nil)
		return
	}

	err := c.ShouldBindJSON(&txrequsets)
	if err != nil || len(txrequsets.Txs) == 0 {
		Response(c, http.StatusBadRequest, fmt.Errorf("request format error: %v", err), nil)
		return
	}
	for i, txReq := range txrequsets.Txs {
		var tx types.Txi
		from, err := common.StringToAddress(txReq.From)
		if err != nil {
			Response(c, http.StatusBadRequest, fmt.Errorf("from address format error: %v", err), nil)
			return
		}

		// empty address
		to := common.BytesToAddress(nil)
		// 0 bigint
		value := math.NewBigInt(0)

		guarantee, ok := math.NewBigIntFromString(txReq.Guarantee, 10)
		if !ok {
			err = fmt.Errorf("convert guarantee to big int error")
		}
		if err != nil {
			Response(c, http.StatusBadRequest, err, nil)
			return
		}

		nonce := txReq.Nonce

		signature := common.FromHex(txReq.Signature)
		if signature == nil {
			Response(c, http.StatusBadRequest, fmt.Errorf("signature format error"), nil)
			return
		}

		pub, err = crypto.Secp256k1PublicKeyFromString(txReq.Pubkey)
		if err != nil {
			Response(c, http.StatusBadRequest, fmt.Errorf("pubkey format error %v", err), nil)
			return
		}

		sig = crypto.SignatureFromBytes(pub.Type, signature)
		if sig.Type != crypto.Signer.GetCryptoType() || pub.Type != crypto.Signer.GetCryptoType() {
			Response(c, http.StatusOK, fmt.Errorf("crypto algorithm mismatch"), nil)
			return
		}
		if !r.SyncerManager.IncrementalSyncer.Enabled {
			Response(c, http.StatusOK, fmt.Errorf("tx is disabled when syncing"), nil)
			return
		}
		tx, err = r.TxCreator.NewTxWithSeal(from, to, value, guarantee, nil, nonce, pub, sig, BaseTokenID)
		if err != nil {
			//try second time
			logrus.WithField("request ", txReq).WithField("tx ", tx).Warn("gen tx failed , try again")
			ok = r.FormatVerifier.VerifySignature(tx)
			if !ok {
				logrus.WithField("request ", txReq).WithField("tx ", tx).Warn("signature invalid")
				Response(c, http.StatusInternalServerError, fmt.Errorf("signature invalid"), nil)
				return
			}
			ok = r.FormatVerifier.VerifySourceAddress(tx)
			if !ok {
				logrus.WithField("request ", txReq).WithField("tx ", tx).Warn("source address invalid")
				Response(c, http.StatusInternalServerError, fmt.Errorf("source address invalid"), nil)
				return
			}
			time.Sleep(time.Microsecond * 2)
			tx, err = r.TxCreator.NewTxWithSeal(from, to, value, guarantee, nil, nonce, pub, sig, BaseTokenID)
			if err != nil {
				logrus.WithField("request ", txReq).WithField("tx ", tx).Warn("gen tx failed")
				Response(c, http.StatusInternalServerError, fmt.Errorf("new tx failed %v", err), nil)
				return
			}
			logrus.WithField("i ", i).WithField("tx", tx).Debugf("tx generated after retry")
			//we don't verify hash , since we calculated the hash
			tx.SetVerified(types.VerifiedFormat)
			r.TxBuffer.ReceivedNewTxChan <- tx
			hashes = append(hashes, tx.GetTxHash())
			continue
		}
		logrus.WithField("i ", i).WithField("tx", tx).Debugf("tx generated")
		//txs = append(txs,tx)
		ok = r.FormatVerifier.VerifySignature(tx)
		if !ok {
			logrus.WithField("request ", txReq).WithField("tx ", tx).Warn("signature invalid")
			Response(c, http.StatusInternalServerError, fmt.Errorf("signature invalid"), nil)
			return
		}
		ok = r.FormatVerifier.VerifySourceAddress(tx)
		if !ok {
			logrus.WithField("request ", txReq).WithField("tx ", tx).Warn("source address invalid")
			Response(c, http.StatusInternalServerError, fmt.Errorf("source address invalid"), nil)
			return
		}
		//we don't verify hash , since we calculated the hash
		tx.SetVerified(types.VerifiedFormat)
		r.TxBuffer.ReceivedNewTxChan <- tx
		hashes = append(hashes, tx.GetTxHash())
	}

	//r.TxBuffer.ReceivedNewTxsChan <- txs

	Response(c, http.StatusOK, nil, hashes)
	return
}
