package main

import (
	"flag"
	"fmt"
	"log"

	"net/http"

	"os"

	"github.com/BurntSushi/toml"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/cli"
	"github.com/versus/gethinx/lib"
	"github.com/versus/gethinx/middle"
	"github.com/versus/gethinx/monitoring"
	"github.com/versus/gethinx/scheduler"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	Version = "v0.1.4"
	Author  = " by Valentyn Nastenko [versus.dev@gmail.com]"
)

var (
	LastBlock      scheduler.EthBlock
	conf           scheduler.Config
	backends       map[string]scheduler.Upstream
	flagConfigFile *string
)

func init() {
	prometheus.MustRegister(monitoring.PromResponse)
	prometheus.MustRegister(monitoring.PromRequest)
	prometheus.MustRegister(monitoring.PromLastBlock)
}

func main() {

	var (
		addr      string
		addrAdmin string
	)

	flagConfigFile = flag.String("c", "./config.toml", "config: path to config file")
	gnrAccKey := flag.Bool("genkey", false, "config: generate access key for agents")
	reloadPtr := flag.Bool("reload", false, "cli: reload only list of servers from config file")
	flag.Parse()

	log.Println("gethinx ", Version, Author)

	if *gnrAccKey {
		fmt.Println("Access Key is ", lib.Key(32))
		os.Exit(0)
	}

	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}

	if *reloadPtr {
		cli.SocketCli(*reloadPtr, &conf)
		os.Exit(0)
	}

	go StartSocketServer()

	if conf.Slack.Use {
		go StartSlackBot()
	}
	if conf.Telegram.Use {
		go StartTelegramBot()
	}
	if govalidator.IsHost(conf.Bind) && govalidator.IsPort(conf.Port) {
		addr = fmt.Sprintf("%s:%s", conf.Bind, conf.Port)
	} else {
		log.Fatalln("Error bind or port in config file")
	}

	if govalidator.IsHost(conf.Bind) && govalidator.IsPort(conf.Port) {
		addrAdmin = fmt.Sprintf("%s:%s", conf.Bind, conf.AdminPort)
	} else {
		log.Fatalln("Error bind or admin port in config file")
	}

	wsAdmin := fmt.Sprintf("ws://%s/ws", addrAdmin)

	initBackendServers()
	GenerateLastBlockAverage()

	go AgentTickerUpstream()

	ar := gin.New()
	ar.LoadHTMLGlob("templates/*")

	ar.GET("/status", func(c *gin.Context) {
		c.HTML(http.StatusOK, "status.tmpl", gin.H{
			"title":     "Gethinx status page",
			"ws_server": wsAdmin,
		})
	})
	ar.GET("/ws", webSocketAdmin)
	ar.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	ar.POST("/api/v1/newblock", setBlock)
	ar.GET("/api/v1/status", getStatus)
	ar.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	go func() {
		err := ar.Run(addrAdmin)
		if err != nil {
			log.Println("Error run admin router: ", err.Error())
		}
	}()

	router := gin.Default()
	router.Use(middle.RequestLogger())
	router.Use(middle.ResponseLogger)
	router.Any("/", reverseProxy)
	//router.POST("/", reverseProxy)

	err := router.Run(addr)
	if err != nil {
		log.Println("Error run gin router: ", err.Error())
	}
}
