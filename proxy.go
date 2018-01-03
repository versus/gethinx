package main

import (
	"github.com/versus/gethinx/lib"
	"github.com/gin-gonic/gin"
	"net/url"
	"net/http/httputil"
	"log"
	"github.com/versus/gethinx/rpc"
	"encoding/json"
)

func reverseProxy(c *gin.Context) {

	var req rpc.JsonRpcMessage

	myreq := lib.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	err := json.Unmarshal(bytes,&req)
	if err != nil {
		log.Println("Error unmarshal  %s", err.Error())
	}

	log.Println(req.Method)
	if req.Method == "eth_getBlockByNumber" {
		hexblock,err := req.GetStringParams(0)
		if err != nil {
			log.Println("Error get Params  %s", err.Error())
		}
		block,err := lib.H2I(hexblock)
		if err != nil {
			log.Println("Error unhex block number  %s", err.Error())
		}
		log.Println("Number  ", block)
	}

	target := "http://127.0.0.1:8080"
	url, err := url.Parse(target)
	if err != nil {
		log.Print("Error parse target %s", err.Error())
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Writer, c.Request)

	}
