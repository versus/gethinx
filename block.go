package main

import (
	"sync"
	"sync/atomic"

	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/lib"
)

type EthBlock struct {
	Dig   int64 `json:"block"`
	Hex   string
	Token string `json:"token"`
	Mutex sync.Mutex
}

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
	//TODO: переделать слайс на мапу с ключем в виде токена
	//Возможно надо создать канал для входящих запросов и увести функцию в горутину
	var agethBlock EthBlock
	myreq := lib.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	if err := json.Unmarshal(bytes, &agethBlock); err != nil {
		log.Println("Error unmarshal ", err.Error())
	}

	srv := backends[agethBlock.Token]
	srv.Mutex.Lock()
	srv.LastBlock = agethBlock.Dig
	srv.HexLastBlock = lib.I2H(agethBlock.Dig)
	srv.Mutex.Unlock()
	backends[agethBlock.Token] = srv

	c.JSON(200, gin.H{
		"blocks": atomic.LoadInt64(&LastBlock.Dig),
	})
}
