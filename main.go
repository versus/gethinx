package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
	"io/ioutil"
	"bytes"
	"io"
	"strconv"
)

func h2i(hex string)  int64 {

	d, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		log.Print("Error parse hex %s", err.Error())
	}
	log.Println(d)
	return d

}

func reverseProxy() gin.HandlerFunc {
	_ = h2i("0xe6")
	return func(c *gin.Context) {
		target := "http://127.0.0.1:8080"
		url, err := url.Parse(target)
		if err != nil {
			log.Print("Error parse target %s", err.Error())
		}
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
		log.Println(readBody(rdr1)) // Print request body
		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	s := buf.String()
	return s
}

func main()  {
	log.Println("hello gethinx")
	router := gin.Default()
	router.Use(RequestLogger())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/", reverseProxy())
	router.Run(":8545")
}
