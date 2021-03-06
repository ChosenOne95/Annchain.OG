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

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func (rpc *RpcController) NewRouter() *gin.Engine {
	router := gin.New()
	if logrus.GetLevel() > logrus.DebugLevel {
		logger := gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: ginLogFormatter,
			Output:    logrus.StandardLogger().Out,
			SkipPaths: []string{"/"},
		})
		router.Use(logger)
	}

	router.Use(gin.RecoveryWithWriter(logrus.StandardLogger().Out))
	return rpc.addRouter(router)
}

func (rpc *RpcController) addRouter(router *gin.Engine) *gin.Engine {
	router.GET("/", rpc.writeListOfEndpoints)
	// init paths here
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("status", rpc.Status)
	router.GET("net_info", rpc.NetInfo)
	router.GET("peers_info", rpc.PeersInfo)
	router.GET("og_peers_info", rpc.OgPeersInfo)
	router.GET("transaction", rpc.Transaction)
	router.GET("transaction_size", rpc.TransactionSize)
	router.GET("confirm", rpc.Confirm)
	router.GET("transactions", rpc.Transactions)
	router.GET("transaction_hashes", rpc.TransactionHashes)
	router.GET("validators", rpc.Validator)
	router.GET("sequencer", rpc.Sequencer)
	router.GET("/v1/sequencer", rpc.SequencerV1)
	router.GET("genesis", rpc.Genesis)
	// broadcast API
	router.POST("new_transaction", rpc.NewTransaction)
	router.GET("new_transaction", rpc.NewTransaction)
	router.POST("new_transactions", rpc.NewTransactions)
	router.POST("new_account", rpc.NewAccount)
	router.POST("new_archive", rpc.NewArchive)
	router.GET("auto_tx", rpc.AutoTx)

	// query API
	router.GET("query", rpc.Query)
	router.GET("query_nonce", rpc.QueryNonce)
	router.GET("query_balance", rpc.QueryBalance)
	router.GET("query_share", rpc.QueryShare)
	router.GET("contract_payload", rpc.ContractPayload)
	router.GET("query_receipt", rpc.QueryReceipt)
	router.POST("query_contract", rpc.QueryContract)
	router.GET("query_contract", rpc.QueryContract)
	router.GET("net_io", rpc.NetIo)

	router.GET("debug", rpc.Debug)
	router.GET("tps", rpc.Tps)
	router.GET("monitor", rpc.Monitor)
	router.GET("sync_status", rpc.SyncStatus)
	router.GET("performance", rpc.Performance)
	router.GET("consensus", rpc.ConStatus)
	router.GET("confirm_status", rpc.ConfirmStatus)

	router.GET("debug/bft_status", rpc.BftStatus)
	router.GET("debug/pool_hashes", rpc.GetPoolHashes)
	router.POST("token/second_offering", rpc.NewSecondOffering) //NewSecondOffering
	router.POST("token/initial_offering", rpc.NewPublicOffering)
	router.POST("token/destroy", rpc.TokenDestroy)
	router.GET("token/latestId", rpc.LatestTokenId)
	router.GET("token/list", rpc.Tokens)
	router.GET("token", rpc.GetToken)
	router.GET("ledger_size", rpc.GetLedgerSize)

	return router

}

// writes a list of available rpc endpoints as an html page
func (rpc *RpcController) writeListOfEndpoints(c *gin.Context) {

	routerMap := map[string]string{
		// info API
		"status":             "",
		"net_info":           "",
		"peers_info":         "",
		"og_peers_info":      "",
		"transaction":        "hash",
		"transaction_size":   "hash",
		"confirm":            "hash",
		"transactions":       "height,address",
		"transaction_hashes": "height",
		"validators":         "",
		"sequencer":          "",
		"/v1/sequencer":      "",
		"genesis":            "",

		"new_transaction": "tx",
		//"new_transaction":  "POSTBODY",
		"new_transactions": "",
		"new_account":      "POSTBODY",
		"new_archive":      "tx",
		"auto_tx":          "interval_us",

		"query":            "query",
		"query_nonce":      "address",
		"query_balance":    "address",
		"query_share":      "pubkey",
		"contract_payload": "payload, abistr",
		"query_receipt":    "hash",
		"query_contract":   "address,data",
		"net_io":           "",
		"debug":            "f",
		"tps":              "",
		"monitor":          "",
		"sync_status":      "",
		"performance":      "",
		"consensus":        "",
		"confirm_status":   "",

		"debug/bft_status":  "",
		"debug/pool_hashes": "",
		"token/latestId":    "",
		"token/list":        "",
		"token":             "id",
		"ledger_size":       "",
	}
	noArgNames := []string{}
	argNames := []string{}
	for name, args := range routerMap {
		if len(args) == 0 {
			noArgNames = append(noArgNames, name)
		} else {
			argNames = append(argNames, name)
		}
	}
	sort.Strings(noArgNames)
	sort.Strings(argNames)
	buf := new(bytes.Buffer)
	buf.WriteString("<html><body>")
	buf.WriteString("<br>Available endpoints:<br>")

	for _, name := range noArgNames {
		link := fmt.Sprintf("http://%s/%s", c.Request.Host, name)
		buf.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a></br>", link, link))
	}

	buf.WriteString("<br>Endpoints that require arguments:<br>")
	for _, name := range argNames {
		link := fmt.Sprintf("http://%s/%s?", c.Request.Host, name)
		args := routerMap[name]
		argNames := strings.Split(args, ",")
		for i, argName := range argNames {
			link += argName + "=_"
			if i < len(argNames)-1 {
				link += "&"
			}
		}
		buf.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a></br>", link, link))
	}
	buf.WriteString("</body></html>")
	//w.Header().Set("Content-Type", "text/html")
	//w.WriteHeader(200)
	//w.Write(buf.Bytes()) // nolint: errcheck
	c.Data(http.StatusOK, "text/html", buf.Bytes())
}
