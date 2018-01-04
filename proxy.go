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

	var req scheduler.JsonRpcMessage

	myreq := lib.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	if err := json.Unmarshal(bytes, &req); err != nil {
		log.Println("Error unmarshal ", err.Error())
	}

	log.Println(req.Method)
	if req.Method == "eth_getBlockByNumber" {
		hexblock, err := req.GetStringParams(0)
		if err != nil {
			log.Println("Error get Params ", err.Error())
		}
		block, err := lib.H2I(hexblock)
		if err != nil {
			log.Println("Error unhex block number ", err.Error())
		}
		log.Println("Number  ", block)
	}

	url, err := target.GetURL()
	if err != nil {
		log.Fatal("Error get URL for ReverseProxy  ", err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Writer, c.Request)

}
