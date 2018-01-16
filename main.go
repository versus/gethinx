package main

import (
	"flag"
	"fmt"
	"log"

	"net/http"

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
	addrAdmin := fmt.Sprintf("%s:%d", conf.Bind, conf.AdminPort)

	initBackendServers()
	GenerateLastBlockAverage()

	go AgentTickerUpstream()

	ar := gin.New()
	ar.LoadHTMLGlob("templates/*")
	ar.GET("/status", func(c *gin.Context) {
		c.HTML(http.StatusOK, "status.tmpl", gin.H{
			"title": "Gethinx status page",
		})
	})
	ar.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	ar.POST("/api/v1/newblock", setBlock)
	ar.GET("/api/v1/status", getStatus)
	go func() {
		err := ar.Run(addrAdmin)
		if err != nil {
			log.Println("Error run admin router: ", err.Error())
		}
	}()

	router := gin.Default()
	router.Use(middle.RequestLogger())
	router.Use(middle.ResponseLogger)

	router.POST("/", reverseProxy)

	err := router.Run(addr)
	if err != nil {
		log.Println("Error run gin router: ", err.Error())
	}
}
