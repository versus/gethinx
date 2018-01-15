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

func generateLastBlockAverage() {
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

func setBlock(c *gin.Context) {
	//TODO: проблема доверия к агенту, возможно надо менять токены  при каждом запросе!!!
	//Возможно надо создать канал для входящих запросов и увести функцию в горутину
	var agethBlock scheduler.EthBlock
	myreq := lib.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	if err := json.Unmarshal(bytes, &agethBlock); err != nil {
		log.Println("Error unmarshal ", err.Error())
	}

	srv := backends[agethBlock.Token]
	srv.Mutex.Lock()
	LastBlock.Mutex.RLock()
	srv.LastBlock = agethBlock.Dig
	srv.HexLastBlock = lib.I2H(agethBlock.Dig)
	srv.TimeUpdate = time.Now().Unix()
	srv.FSM.Event("up")
	LastBlock.Mutex.RUnlock()
	srv.Mutex.Unlock()
	backends[agethBlock.Token] = srv
	generateLastBlockAverage()

	c.JSON(200, gin.H{
		"average blocks": atomic.LoadInt64(&LastBlock.Dig),
	})
}
