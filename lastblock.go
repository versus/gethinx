package main

import (
	"context"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/versus/gethinx/lib"
	"github.com/versus/gethinx/monitoring"
	"github.com/versus/gethinx/scheduler"
)

//GenerateLastBlockAverage is func for generate average ariphemic number block
func GenerateLastBlockAverage() {
	sum, count := int64(0), int64(0)
	average := int64(0)

	for _, srv := range backends {
		if srv.FSM.Current() == "active" {
			sum = sum + srv.LastBlock
			count++
		}
	}
	if count != 0 {
		average = int64(sum / count)
	}
	LastBlock.Mutex.Lock()
	LastBlock.Dig = average
	LastBlock.Hex = lib.I2H(average)
	LastBlock.TimeUpdate = time.Now().Unix()
	monitoring.PromLastBlock.Set(float64(LastBlock.Dig))
	LastBlock.Mutex.Unlock()

}

func AgentTickerUpstream() {
	log.Println(conf.Ticker)
	tick := time.Tick(time.Second * time.Duration(conf.Ticker))
	for {
		select {
		case <-tick:
			for key, srv := range backends {
				if srv.FSM.Current() == "active" {
					lastTimeUpdate := time.Unix(srv.TimeUpdate, 0)
					now := time.Now()
					diff := now.Sub(lastTimeUpdate)
					//log.Println(srv.Target, "time sub is ", int64(diff/1000000000), "suspend time is ", int64(conf.Suspend))
					if int64(diff/1000000000) > int64(conf.Suspend) {
						srv.Mutex.Lock()
						srv.FSM.Event("suspend")
						srv.Mutex.Unlock()
					}
					backends[key] = srv
				}
			}
			checkAlive()

		}
	}
}

func checkAlive() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for key, srv := range backends {
		srv.GetTargetLastBlock(ctx, &LastBlock)
		//log.Println("checkAlive: ", srv.Target, " is ", srv.FSM.Current())
		backends[key] = srv
	}
	GenerateLastBlockAverage()
}

func ReloadBackendServers(configFile *string) {
	if _, err := toml.DecodeFile(*configFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}
	log.Println("add server  with ", len(conf.Servers))
	for k := range backends {
		delete(backends, k)
	}
	generatorBackend()

}

func initBackendServers() {
	if len(conf.Servers) == 0 {
		log.Fatalln("Servers for backend is not defined")
	}

	backends = make(map[string]scheduler.Upstream, len(conf.Servers))
	generatorBackend()

}

func generatorBackend() {
	for _, srvValue := range conf.Servers {
		backends[srvValue.Token] = *scheduler.NewUpstream(srvValue.IP, srvValue.Port, srvValue.Weight, srvValue.Token, srvValue.Hostname)
		log.Println("add server  with ", backends[srvValue.Token].Target)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		target := backends[srvValue.Token]
		target.GetTargetLastBlock(ctx, &LastBlock)
		backends[srvValue.Token] = target
	}

}
