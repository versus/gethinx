package main

import (
	"log"
	"sync/atomic"

	"flag"

	"fmt"

	"context"

	"time"

	"sync"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/lib"
	"github.com/versus/gethinx/middle"
	"github.com/versus/gethinx/scheduler"
)

const gethinxVersion = "0.0.1"

var (
	LastBlock EthLastBlock
	conf      scheduler.Config
	backends  map[int]scheduler.Upstream
)

type EthLastBlock struct {
	Dig   int64
	Hex   string
	Mutex sync.Mutex
}

func setBlock(c *gin.Context) {
	//TODO: распарсить последний блок и занести в мапу серверов бекенда
	//TODO: произвести расчет нового среднего блока
	//TODO: проблема доверия к агенту, возможно надо менять токены  прик аждом запросе!!!
	c.JSON(200, gin.H{
		"blocks": atomic.LoadInt64(&LastBlock.Dig),
	})
}

func generateLastBlockAverage() {
	sum, count := int64(0), int64(0)
	average := int64(0)

	for _, srv := range backends {
		log.Println(" fsm is ", srv.FSM.Current())
		if srv.FSM.Current() == "active" {
			sum = sum + srv.LastBlock
			count++
		}
	}
	log.Println("sum is ", sum, " count is ", count)
	if count != 0 {
		average = int64(sum / count)
	}
	LastBlock.Mutex.Lock()
	LastBlock.Dig = average
	LastBlock.Hex = lib.I2H(average)
	LastBlock.Mutex.Unlock()

}

func initBackendServers() {
	if len(conf.Servers) == 0 {
		log.Fatalln("Servers for backend is not defined")
	}

	backends = make(map[int]scheduler.Upstream, len(conf.Servers))
	srvKey := 0
	for _, srvValue := range conf.Servers {
		backends[srvKey] = *scheduler.NewUpstream(srvValue.IP, srvValue.Port, srvValue.Weight, srvValue.Token)
		log.Println("add server ", srvKey, " with ", backends[srvKey].Target)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		target := backends[srvKey]
		target.GetTargetLastBlock(ctx)
		backends[srvKey] = target
		srvKey++
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
