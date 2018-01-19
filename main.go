package main

import (
	"flag"
	"fmt"
	"log"

	"net/http"

	"net"

	"github.com/BurntSushi/toml"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/middle"
	"github.com/versus/gethinx/scheduler"
	"os"
	"syscall"
	"os/signal"
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
	var (
		addr      string
		addrAdmin string
	)

	flagConfigFile := flag.String("c", "./config.toml", "config: path to config file")
	flag.Parse()

	log.Println("gethinx ", Version, Author)

	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}

	ln, err := net.Listen("unix", "/tmp/gethinx.sock")
	if err != nil {
		log.Fatal("Listen error: /tmp/gethinx.sock ", err.Error())
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
	}(ln, sigc)

	go StartSocketServer(ln)

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

	err = router.Run(addr)
	if err != nil {
		log.Println("Error run gin router: ", err.Error())
	}
}
