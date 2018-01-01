package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
	"github.com/versus/gethinx/lib"
	"github.com/versus/gethinx/middle"
	"github.com/versus/gethinx/rpc"
	"encoding/json"
)

var (
	numBlocks int64 = 3644
)




func reverseProxy() gin.HandlerFunc {
	_ = lib.H2I("0xe6")
	_ = lib.I2H(230)

	return func(c *gin.Context) {
		target := "http://127.0.0.1:8080"
		url, err := url.Parse(target)
		if err != nil {
			log.Print("Error parse target %s", err.Error())
		}

        var req rpc.JsonRpcMessage
		js := `{"jsonrpc":"2.0","id":6,"method":"eth_getBlockByNumber","params":["0xe3c",false]}`
		bytes := []byte(js)
		err = json.Unmarshal(bytes,&req)
		if err != nil {
			log.Println("Error unmarshal  %s", err.Error())
		}
		log.Println(req.Method)
		log.Println(string(req.Params[0]))

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main()  {
	log.Println("hello gethinx")
	router := gin.Default()
	router.Use(middle.RequestLogger())
	router.Use(middle.ResponseLogger)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/", reverseProxy())
	router.Run(":8545")
}

