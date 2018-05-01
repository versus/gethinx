package backend

import (
	"context"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/buffer"
	"github.com/versus/gethinx/config"
	"github.com/versus/gethinx/ethblock"
	"log"
	"time"
)

type BackendList map[string]Upstream

func NewBackendServers(count int) *BackendList {
	if count == 0 {
		log.Fatalln("Servers for backend is not defined")
	}

	backends := make(BackendList, count)
	return &backends

}

func (b BackendList) GeneratorBackend(conf config.Config) {
	for _, srvValue := range conf.Servers {
		b[srvValue.Token] = *NewUpstream(srvValue.IP, srvValue.Port, srvValue.Weight, srvValue.Token, srvValue.Hostname)
		log.Println("add server  with ", b[srvValue.Token].Target)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		target := b[srvValue.Token]
		target.TargetLastBlock(ctx, &ethblock.LastBlock)
		b[srvValue.Token] = target
	}

}

func (b BackendList) ReloadBackendServers(configFile string) (config.Config, error) {
	var conf config.Config
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		return conf, err
	}
	log.Println("add server  with ", len(conf.Servers))
	for k := range b {
		delete(b, k)
	}
	b.GeneratorBackend(conf)
	return conf, nil
}

func SetBlock(c *gin.Context) {
	//TODO: проблема доверия к агенту, возможно надо менять токены  при каждом запросе!!!
	var agethBlock ethblock.EthBlock
	myreq := buffer.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	if err := json.Unmarshal(bytes, &agethBlock); err != nil {
		log.Println("Error unmarshal ", err.Error())
	}
	/*
		srv := backends[agethBlock.Token]
		LastBlock.Mutex.RLock()
		srv.LastBlock = agethBlock.Dig
		srv.HexLastBlock = I2H(agethBlock.Dig)
		srv.TimeUpdate = time.Now().Unix()
		srv.FSM.Event("up")
		LastBlock.Mutex.RUnlock()
		backends[agethBlock.Token] = srv
		GenerateLastBlockAverage()

		c.JSON(200, gin.H{
			"average blocks": atomic.LoadInt64(&LastBlock.Dig),
		})
	*/
}
