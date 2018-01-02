package main

import (
	"github.com/versus/gethinx/lib"
	"github.com/gin-gonic/gin"
	"net/url"
	"net/http/httputil"
	"log"
	"github.com/versus/gethinx/rpc"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

func reverseProxy(c *gin.Context) {

	var req rpc.JsonRpcMessage
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Print("Error parse c.Request.Body  %s", err.Error())
	}
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	reqBody := lib.ReadBody(rdr1)
	log.Println(reqBody) // Print request body
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
	c.Request.Body = rdr2

	_ = lib.H2I("0xe6")
	_ = lib.I2H(230)

	target := "http://127.0.0.1:8080"
	url, err := url.Parse(target)
	if err != nil {
		log.Print("Error parse target %s", err.Error())
	}

	bytes := []byte(reqBody)
	err = json.Unmarshal(bytes,&req)
	if err != nil {
		log.Println("Error unmarshal  %s", err.Error())
	}


	log.Println(req.Method)
	if req.Method == "eth_getBlockByNumber" {
		log.Println(string(req.Params[0]))
	}


	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Writer, c.Request)

	}
