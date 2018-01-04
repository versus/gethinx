package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/lib"
	"github.com/versus/gethinx/scheduler"
	"log"
	"net/http/httputil"
	"net/url"
)

func getUrl(req lib.MyRequestBody) (*url.URL, error) {
	_ = req.Body
	target := "http://127.0.0.1:8080"
	url, err := url.Parse(target)
	return url, err
}

func reverseProxy(c *gin.Context) {

	var req scheduler.JsonRpcMessage

	myreq := lib.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	err := json.Unmarshal(bytes, &req)
	if err != nil {
		log.Println("Error unmarshal  %s", err.Error())
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

	url, err := getUrl(myreq)
	if err != nil {
		log.Fatal("Error get URL for ReverseProxy  ", err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Writer, c.Request)

}
