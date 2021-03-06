package main

import (
	"encoding/json"
	"log"
	"net/http/httputil"

	"net/url"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/lib"
	"github.com/versus/gethinx/scheduler"
)

func reverseProxy(c *gin.Context) {

	var req scheduler.JSONRPCMessage
	var block int64 = -1

	var url *url.URL
	var err error

	myreq := lib.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	if err := json.Unmarshal(bytes, &req); err != nil {
		log.Println("Error unmarshal ", err.Error())
	}

	log.Println(req.Method)

	LastBlock.Mutex.RLock()
	hex := LastBlock.Hex
	LastBlock.Mutex.RUnlock()

	if req.Method == "eth_blockNumber" {
		c.JSON(200, gin.H{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  hex,
		})
		return
	}

	if req.Method == "eth_getBlockByNumber" {
		// Req:  {"jsonrpc":"2.0","id":13,"method":"eth_getBlockByNumber","params":[null,false]}
		// Req:  {"jsonrpc":"2.0","id":14,"method":"eth_getBlockByNumber","params":["latest",false]}

		hexblock, err := req.GetStringParams(0)
		if err != nil {
			log.Println("Error get Params ", err.Error())
		}
		if hexblock != "latest" {
			block, err = lib.H2I(hexblock)
			if err != nil {
				log.Println("Error unhex block number ", err.Error())
			}
		}

	}
	url, err = scheduler.GetTargetNode(backends, block, &LastBlock, conf.MaxResponseTime)
	if err != nil {
		log.Println("Error get URL for ReverseProxy  ", err.Error())
		c.Writer.WriteHeader(http.StatusBadGateway)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Writer, c.Request)

}
