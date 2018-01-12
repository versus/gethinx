package main

import (
	"sync/atomic"

	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

func getStatus(c *gin.Context) {
	var backendsServerList []string
	for _, server := range backends {
		server.Mutex.Lock()
		server.RealState = server.FSM.Current()
		server.Mutex.Unlock()
		jsrv, err := json.Marshal(server)
		if err != nil {
			log.Print("Error marchal backendsServerList for status ", err.Error())
		}
		backendsServerList = append(backendsServerList, string(jsrv))
	}
	c.JSON(200, gin.H{
		"ServerList":       backendsServerList,
		"LastBlockAverage": atomic.LoadInt64(&lastBlock),
	})
}
