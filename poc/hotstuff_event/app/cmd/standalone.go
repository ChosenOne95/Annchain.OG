package cmd

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/annchain/OG/poc/hotstuff_event"
	core "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func MakeStandalonePartner(myId int, N int, F int, hub hotstuff_event.Hub, peerIds []string) *hotstuff_event.Partner {
	logger := hotstuff_event.SetupOrderedLog(myId)
	ledger := &hotstuff_event.Ledger{
		Logger: logger,
	}
	ledger.InitDefault()

	safety := &hotstuff_event.Safety{
		Ledger: ledger,
		Logger: logger,
	}

	blockTree := &hotstuff_event.BlockTree{
		Ledger:    ledger,
		F:         F,
		Logger:    logger,
		MyIdIndex: myId,
	}
	blockTree.InitDefault()
	blockTree.InitGenesisOrLatest()

	proposerElection := &hotstuff_event.ProposerElection{N: N}

	paceMaker := &hotstuff_event.PaceMaker{
		PeerIds:          peerIds,
		MyIdIndex:        myId,
		CurrentRound:     1, // must be 1 which is AFTER GENESIS
		Safety:           safety,
		MessageHub:       hub,
		BlockTree:        blockTree,
		ProposerElection: proposerElection,
		Logger:           logger,
		Partner:          nil, // fill later
	}
	paceMaker.InitDefault()

	blockTree.PaceMaker = paceMaker

	partner := &hotstuff_event.Partner{
		PeerIds:          peerIds,
		MessageHub:       hub,
		Ledger:           ledger,
		MyIdIndex:        myId,
		N:                N,
		F:                F,
		PaceMaker:        paceMaker,
		Safety:           safety,
		BlockTree:        blockTree,
		ProposerElection: proposerElection,
		Logger:           logger,
	}
	partner.InitDefault()

	paceMaker.Partner = partner
	safety.Partner = partner

	return partner
}

// runCmd represents the run command
var standaloneCmd = &cobra.Command{
	Use:   "standalone",
	Short: "Start a standalone node",
	Long:  `Start a standalone node and communicate with other standalone nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLogger()

		peers := readList(viper.GetString("list"))
		//total := len(peers)

		priv, id := loadPrivateKey()

		p2p := &hotstuff_event.PhysicalCommunicator{
			Port:       viper.GetInt("port"),
			PrivateKey: priv,
		}
		p2p.InitDefault()

		hub := &hotstuff_event.LogicalCommunicator{
			PhysicalCommunicator: p2p,
			MyId:                 id,
		}

		hub.InitDefault()
		hub.Start()

		// init me before init peers
		hub.PhysicalCommunicator.Start()

		go func() {
			for {
				// now broadcast constantly
				hub.Broadcast(&hotstuff_event.Msg{
					Typev:    hotstuff_event.String,
					Sig:      hotstuff_event.Signature{},
					SenderId: hub.MyId,
					Content: &hotstuff_event.ContentString{
						Content: fmt.Sprintf("MSG %s->%s", hub.MyId, time.Now().String())},
				}, "")
				time.Sleep(time.Second * 2)
			}
		}()
		go func() {
			// preconnect peers
			for _, peer := range peers {
				hub.PhysicalCommunicator.SuggestConnection(peer)
			}
		}()

		go func() {
			messageChannel, _ := hub.GetChannel(hub.MyId)
			for {
				v := <-messageChannel
				fmt.Println("I received " + v.Content.String())
			}
		}()

		//partners := make([]*hotstuff_event.Partner, total)

		//partner := MakeStandalonePartner(viper.GetInt("mei"), total, total/3, hub)
		//go partner.Start()

		// prevent sudden stop. Do your clean up here
		var gracefulStop = make(chan os.Signal)

		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)

		func() {
			sig := <-gracefulStop
			log.Warnf("caught sig: %+v", sig)
			log.Warn("Exiting... Please do no kill me")
			//for _, partner := range partners {
			//	partner.Stop()
			//}
			os.Exit(0)
		}()

	},
}

func readList(filename string) (peers []string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		peers = append(peers, line)

	}
	return

}

func loadPrivateKey() (core.PrivKey, string) {
	// read key file
	keyFile := viper.GetString("file")
	bytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}

	pi := &hotstuff_event.PrivateInfo{}
	err = json.Unmarshal(bytes, pi)
	if err != nil {
		panic(err)
	}

	privb, err := hex.DecodeString(pi.PrivateKey)
	if err != nil {
		panic(err)
	}

	priv, err := core.UnmarshalPrivateKey(privb)
	if err != nil {
		panic(err)
	}
	return priv, pi.Id
}

func init() {
	rootCmd.AddCommand(standaloneCmd)
	standaloneCmd.Flags().StringP("list", "l", "peers.lst", "Partners to be started in file list")
	_ = viper.BindPFlag("list", standaloneCmd.Flags().Lookup("list"))

	standaloneCmd.Flags().StringP("file", "f", "id.key", "my key file")
	_ = viper.BindPFlag("file", standaloneCmd.Flags().Lookup("file"))

	standaloneCmd.Flags().IntP("port", "p", 3301, "Local IO port")
	_ = viper.BindPFlag("port", standaloneCmd.Flags().Lookup("port"))
}
