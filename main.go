package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/middle"
	"github.com/versus/gethinx/scheduler"
)

const (
	Version = "v0.0.1"
	Author  = " by Valentyn Nastenko [versus.dev@gmail.com]"
)

var (
	LastBlock scheduler.EthBlock
	conf      scheduler.Config
	backends  map[string]scheduler.Upstream
)

func initBackendServers() {
	if len(conf.Servers) == 0 {
		log.Fatalln("Servers for backend is not defined")
	}

	backends = make(map[string]scheduler.Upstream, len(conf.Servers))
	for _, srvValue := range conf.Servers {
		backends[srvValue.Token] = *scheduler.NewUpstream(srvValue.IP, srvValue.Port, srvValue.Weight, srvValue.Token)
		log.Println("add server  with ", backends[srvValue.Token].Target)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		target := backends[srvValue.Token]
		target.GetTargetLastBlock(ctx)
		backends[srvValue.Token] = target
	}
}

func main() {
	//TODO: create flag to reload config only
	//TODO: флаг работы без агента (замедление работы с новыми блоками!!!!)
	//TODO: create socket for client command
	flagConfigFile := flag.String("c", "./config.toml", "config: path to config file")
	flag.Parse()

	log.Println("gethinx ", Version, Author)

	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}

	addr := fmt.Sprintf("%s:%d", conf.Bind, conf.Port)

	initBackendServers()
	generateLastBlockAverage()

	router := gin.Default()
	router.Use(middle.RequestLogger())
	router.Use(middle.ResponseLogger)
	router.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/api/v1/newblock", setBlock)
	router.POST("/", reverseProxy)
	router.GET("/api/v1/status", getStatus)
	err := router.Run(addr)
	if err != nil {
		log.Println("Error run gin router: ", err.Error())
	}
}
