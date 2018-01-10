package main

import (
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"flag"

	"fmt"

	"context"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/middle"
	"github.com/versus/gethinx/scheduler"
)

const gethinxVersion = "0.0.1"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("gethinx ", gethinxVersion, " (c)2018 Valentyn Nastenko")

	var conf scheduler.Config
	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}

	addr := fmt.Sprintf("%s:%d", conf.Bind, conf.Port)


	//TODO: вынести в инициализацию сервера бекенда
	target = scheduler.NewUpstream(conf.Servers["alpha"].IP, conf.Servers["alpha"].Port, conf.Servers["alpha"].Weight)
	target.GetTargetLastBlock(ctx)
	log.Println("target state is ", target.FSM.Current())
	log.Println("target last block ", target.LastBlock.String())

	router := gin.Default()
	router.Use(middle.RequestLogger())
	router.Use(middle.ResponseLogger)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/api/v1/newblock", setBlock)
	router.POST("/", reverseProxy)
	router.GET("/status", getStatus)
	err := router.Run(addr)
	if err != nil {
		log.Println("Error run gin router: ", err.Error())
	}
}
