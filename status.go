package main

import (
	"sync/atomic"

	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

func JsonBackends() []string {
	var backendsServerList []string
	for _, server := range backends {
		jsrv, err := json.Marshal(server)
		if err != nil {
			log.Print("Error marchal backendsServerList for status ", err.Error())
		}
		backendsServerList = append(backendsServerList, string(jsrv))
	}
	return backendsServerList
}

func getStatus(c *gin.Context) {

	backendsServerList := JsonBackends()
	c.JSON(200, gin.H{
		"ServerList":       backendsServerList,
		"LastBlockAverage": atomic.LoadInt64(&LastBlock.Dig),
	})
}
