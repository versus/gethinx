package main

import (
	"log"
	"sync/atomic"

	"flag"

	"fmt"

	"context"

	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/middle"
	"github.com/versus/gethinx/scheduler"
)

const gethinxVersion = "0.0.1"

var (
	numBlocks int64 = 3644
	target    *scheduler.Upstream
	conf      scheduler.Config
	backends  map[string]scheduler.Upstream
)

func setBlock(c *gin.Context) {
	//TODO: распарсить последний блок и занести в мапу серверов бекенда
	//TODO: произвести расчет нового среднего блока
	//TODO: проблема доверия к агенту!!!
	c.JSON(200, gin.H{
		"blocks": atomic.LoadInt64(&numBlocks),
	})
}

func initBackendServers() {
	if len(conf.Servers) == 0 {
		log.Fatalln("Servers for backend is not defined")
	}
	backends = make(map[string]scheduler.Upstream, len(conf.Servers))
	for srvKey, srvValue := range conf.Servers {
		backends[srvKey] = *scheduler.NewUpstream(srvValue.IP, srvValue.Port, srvValue.Weight)
		log.Println("add server ", srvKey, " with ", backends[srvKey].Target)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		target := backends[srvKey]
		target.GetTargetLastBlock(ctx)
		log.Println("target ", target.Target, " state is ", target.FSM.Current())
		log.Println("target ", target.Target, "last block ", target.LastBlock)
	}
}

func main() {
	//TODO: create flag to reload config only
	//TODO: флаг работы без агента (замедление работы с новыми блоками!!!!)
	//TODO: create socket for client command
	flagConfigFile := flag.String("c", "./config.toml", "config: path to config file")
	flag.Parse()

	log.Println("gethinx ", gethinxVersion, " (c)2018 Valentyn Nastenko")

	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}

	addr := fmt.Sprintf("%s:%d", conf.Bind, conf.Port)

	log.Println("count of servers: ", len(conf.Servers))
	initBackendServers()

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
