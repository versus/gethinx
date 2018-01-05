package main

import (
	"log"
	"strconv"
	"sync/atomic"

	"flag"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/middle"
	"github.com/versus/gethinx/scheduler"
)

var numBlocks int64 = 3644
var target *scheduler.Upstream

func setBlock(c *gin.Context) {
	newBlocks, err := strconv.ParseInt(c.PostForm("block"), 0, 64)
	atomic.StoreInt64(&numBlocks, newBlocks)
	if err != nil {
		log.Println("Error parse post block  ", err.Error())
	}
	c.JSON(200, gin.H{
		"blocks": atomic.LoadInt64(&numBlocks),
	})
}

func main() {
	flagConfigFile := flag.String("c", "./config.toml", "config: path to config file")
	flag.Parse()

	log.Println("gethinx v0.0.1 (c)2018 Valentyn Nastenko")

	var conf scheduler.Config
	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}
	log.Println("Toml cats: ", conf.Cats)

	target = scheduler.NewUpstream("127.0.0.1", "8080", "1")
	log.Println("target state is ", target.FSM.Current())

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
	err := router.Run(":8545")
	if err != nil {
		log.Println("Error run gin router: ", err.Error())
	}
}
