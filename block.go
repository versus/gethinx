package main

import (
	"sync/atomic"

	"encoding/json"
	"log"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/lib"
	"github.com/versus/gethinx/scheduler"
)

func setBlock(c *gin.Context) {
	//TODO: проблема доверия к агенту, возможно надо менять токены  при каждом запросе!!!
	var agethBlock scheduler.EthBlock
	myreq := lib.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	if err := json.Unmarshal(bytes, &agethBlock); err != nil {
		log.Println("Error unmarshal ", err.Error())
	}

	srv := backends[agethBlock.Token]
	LastBlock.Mutex.RLock()
	srv.LastBlock = agethBlock.Dig
	srv.HexLastBlock = lib.I2H(agethBlock.Dig)
	srv.TimeUpdate = time.Now().Unix()
	srv.FSM.Event("up")
	LastBlock.Mutex.RUnlock()
	backends[agethBlock.Token] = srv
	GenerateLastBlockAverage()

	c.JSON(200, gin.H{
		"average blocks": atomic.LoadInt64(&LastBlock.Dig),
	})

}
