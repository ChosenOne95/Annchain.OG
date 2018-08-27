package node

import (
	"github.com/annchain/OG/og"
	"github.com/annchain/OG/p2p"
	"github.com/annchain/OG/rpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Node is the basic entrypoint for all modules to start.
type Node struct {
	og        *og.Og
	p2pServer *p2p.Server
	rpcServer *rpc.RpcServer
}

func NewNode() *Node {
	n := new(Node)
	if viper.GetBool("p2p.enabled") {
		privKey := getNodePrivKey()
		n.p2pServer = NewP2PServer(privKey)
	}
	if viper.GetBool("rpc.enabled") {
		n.rpcServer = rpc.NewRpcServer(viper.GetString("rpc.port"))
	}
	n.og = new(og.Og)
	return n
}

func (n *Node) Start() {
	// start myself first and then provide service and do p2p
	n.og.Start()

	if n.p2pServer != nil {
		logrus.Info("Starting p2p server")
		n.p2pServer.Start()
	}
	if n.rpcServer != nil {
		logrus.Info("Starting rpc server")
		n.rpcServer.Start()
	}
	logrus.Info("Node Started")
}
func (n *Node) Stop() {
	if n.rpcServer != nil {
		n.rpcServer.Stop()
	}
	if n.p2pServer != nil {
		n.p2pServer.Stop()
	}
	n.og.Stop()
	logrus.Info("Node Stopped")
}
