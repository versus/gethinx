package main

import (
	"context"
	"log"
	"time"

	"github.com/versus/gethinx/lib"
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
	LastBlock.Mutex.Unlock()
}

func AgentTickerUpstream() {
	tick := time.Tick(time.Second * 10)
	for {
		select {
		case <-tick:
			alive := 0
			for key, srv := range backends {
				if srv.FSM.Current() == "active" {
					alive++
					lastTimeUpdate := time.Unix(srv.TimeUpdate, 0)
					now := time.Now()
					diff := now.Sub(lastTimeUpdate)
					log.Println(srv.Target, "time sub is ", int64(diff/1000000000), "suspend time is ", int64(conf.Suspend))
					if int64(diff/1000000000) > int64(conf.Suspend) {
						srv.Mutex.Lock()
						srv.FSM.Event("suspend")
						srv.Mutex.Unlock()
					}
					backends[key] = srv
				}
			}
			if alive == 0 {
				checkAlive()
			}

		}
	}
}

func checkAlive() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for key, srv := range backends {
		srv.GetTargetLastBlock(ctx)
		log.Println("checkAlive: ", srv.Target, " is ", srv.FSM.Current())
		backends[key] = srv
	}
	GenerateLastBlockAverage()
}

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
