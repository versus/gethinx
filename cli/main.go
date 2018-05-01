package main

import (
	"flag"
	"fmt"
	"log"

	"os"

	"github.com/BurntSushi/toml"
	"github.com/asaskevich/govalidator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/versus/gethinx"
	"github.com/versus/gethinx/backend"
	"github.com/versus/gethinx/config"
	"github.com/versus/gethinx/monitoring"
	"github.com/versus/gethinx/proxyserver"
)

const (
	Version = "v0.2.0"
	Author  = " by Valentyn Nastenko [versus.dev@gmail.com]"
)

var (
	conf           config.Config
	backends       *backend.BackendList
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
		fmt.Println("Access Key is ", gethinx.Key(32))
		os.Exit(0)
	}

	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}

	if *reloadPtr {
		gethinx.SocketCli(*reloadPtr, &conf)
		os.Exit(0)
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

	go gethinx.StartSocketServer(*flagConfigFile, conf.SocketPath)

	backends = backend.NewBackendServers(len(conf.Servers))
	proxyserver.GenerateLastBlockAverage()

	go AgentTickerUpstream()

	StartApi(addr, addrAdmin)
}
