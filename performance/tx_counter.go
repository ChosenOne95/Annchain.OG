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
package performance

import (
	"github.com/annchain/OG/common/goroutine"
	"github.com/annchain/OG/types"
	"go.uber.org/atomic"
	"time"
)

type TxCounter struct {
	TxGenerated        atomic.Uint64 `json:"txGenerated"`        // tx generated by me
	TxReceived         atomic.Uint64 `json:"txReceived"`         // tx generated by me and others
	TxConfirmed        atomic.Uint64 `json:"txConfirmed"`        // tx confirmed
	SequencerGenerated atomic.Uint64 `json:"sequencerGenerated"` // sequencer generated by me
	SequencerReceived  atomic.Uint64 `json:"sequencerReceived"`  // sequencer generated by me and others
	SequencerConfirmed atomic.Uint64 `json:"sequencerConfirmed"` // sequencer confirmed
	StartupTime        time.Time     `json:"startupTime"`        // timestamp of the program started.

	// listeners
	NewTxReceivedChan  chan types.Txi   `json:"-"`
	NewTxGeneratedChan chan types.Txi   `json:"-"`
	BatchConfirmedChan chan []types.Txi `json:"-"`
	quit               chan bool        `json:"-"`
}

func NewTxCounter() *TxCounter {
	return &TxCounter{
		BatchConfirmedChan: make(chan []types.Txi),
		NewTxReceivedChan:  make(chan types.Txi),
		NewTxGeneratedChan: make(chan types.Txi),
		quit:               make(chan bool),
		StartupTime:        time.Now(),
	}
}

func (t *TxCounter) loop() {
	for {
		select {
		case <-t.quit:
			break
		case tx := <-t.NewTxReceivedChan:
			switch tx.GetType() {

			case types.TxBaseTypeSequencer:
				t.SequencerReceived.Inc()
			default:
				t.TxReceived.Inc()
				//logrus.WithField("type", tx.GetType()).Debug("Unknown tx type")
			}
		case batch := <-t.BatchConfirmedChan:
			for _, tx := range batch {
				switch tx.GetType() {
				case types.TxBaseTypeSequencer:
					t.SequencerConfirmed.Inc()
				default:
					t.TxConfirmed.Inc()
					//logrus.WithField("type", tx.GetType()).Debug("Unknown tx type")
				}
			}
		case tx := <-t.NewTxGeneratedChan:
			switch tx.GetType() {
			case types.TxBaseTypeSequencer:
				t.SequencerGenerated.Inc()
			default:
				t.TxGenerated.Inc()
				//logrus.WithField("type", tx.GetType()).Debug("Unknown tx type")
			}
		}
	}
}

func (t *TxCounter) Start() {
	goroutine.New(t.loop)
}

func (t *TxCounter) Stop() {
	close(t.quit)
}

func (*TxCounter) Name() string {
	return "TxCounter"
}

func (t *TxCounter) GetBenchmarks() map[string]interface{} {
	return map[string]interface{}{
		"txGenerated":        t.TxGenerated.Load(),
		"txReceived":         t.TxReceived.Load(),
		"txConfirmed":        t.TxConfirmed.Load(),
		"sequencerGenerated": t.SequencerGenerated.Load(),
		"sequencerReceived":  t.SequencerReceived.Load(),
		"sequencerConfirmed": t.SequencerConfirmed.Load(),
		"startupTime":        t.StartupTime.Unix(),
	}
}
