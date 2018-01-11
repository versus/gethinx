package main

import (
	"encoding/json"
	"log"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/lib"
	"github.com/versus/gethinx/scheduler"
)

func reverseProxy(c *gin.Context) {

	var req scheduler.JSONRPCMessage
	var block int64

	myreq := lib.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	if err := json.Unmarshal(bytes, &req); err != nil {
		log.Println("Error unmarshal ", err.Error())
	}

	log.Println(req.Method)
	if req.Method == "eth_blockNumber" {
		c.JSON(200, gin.H{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  "0x82a",
		})
		return
	}
	if req.Method == "eth_getBlockByNumber" {
		//TODO check null and latest
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
		log.Println("Number  ", block)
	}

	//TODO setup our target url
	target := backends[0]
	log.Println("Target host: ", target.Target)
	url, err := target.GetURL()
	if err != nil {
		log.Fatal("Error get URL for ReverseProxy  ", err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Writer, c.Request)

}
