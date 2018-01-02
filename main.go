package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/middle"
	"strconv"
	"sync/atomic"
)

var numBlocks int64 = 3644


func setBlock(c *gin.Context){
	newBlocks,err := strconv.ParseInt(c.PostForm("block"), 0, 64)
	atomic.StoreInt64(&numBlocks,newBlocks)
	if err != nil {
		log.Println("Error parse post block  %s", err.Error())
	}
	c.JSON(200, gin.H{
			"blocks": atomic.LoadInt64(&numBlocks),
		})
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
	router.POST("/block", setBlock)
	router.POST("/", reverseProxy)
	router.GET("/status", getStatus)
	router.Run(":8545")
}

