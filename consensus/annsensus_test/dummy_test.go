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
package annsensus_test

import (
	"github.com/annchain/OG/common"
	"github.com/annchain/OG/common/math"
	"github.com/annchain/OG/og/protocol/ogmessage"
)

type DummyDag struct {
}

func (d *DummyDag) GetTx(hash common.Hash) ogmessage.Txi {
	return nil
}

func (d *DummyDag) GetTxByNonce(addr common.Address, nonce uint64) ogmessage.Txi {
	return nil
}

func (d *DummyDag) GetLatestNonce(addr common.Address) (uint64, error) {
	return 0, nil
}

func (d *DummyDag) GetSequencerByHeight(id uint64) *ogmessage.Sequencer {
	return &ogmessage.Sequencer{
		TxBase: ogmessage.TxBase{Height: id},
	}
}

func (d *DummyDag) GetTxisByNumber(id uint64) ogmessage.Txis {
	var txis ogmessage.Txis
	txis = append(txis, ogmessage.RandomTx(), ogmessage.RandomTx())
	return txis
}

func (d *DummyDag) LatestSequencer() *ogmessage.Sequencer {
	return ogmessage.RandomSequencer()
}

func (d *DummyDag) GetSequencer(hash common.Hash, id uint64) *ogmessage.Sequencer {
	return &ogmessage.Sequencer{
		TxBase: ogmessage.TxBase{Height: id,
			Hash: hash},
	}
}

func (d *DummyDag) Genesis() *ogmessage.Sequencer {
	return &ogmessage.Sequencer{
		TxBase: ogmessage.TxBase{Height: 0},
	}
}

func (d *DummyDag) GetHeight() uint64 {
	return 0
}

func (d *DummyDag) GetSequencerByHash(hash common.Hash) *ogmessage.Sequencer {
	return nil
}

func (d *DummyDag) GetBalance(addr common.Address, tokenId int32) *math.BigInt {
	return math.NewBigInt(100000)
}
