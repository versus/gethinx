package main

import (
	"sync/atomic"

	"encoding/json"
	"log"

	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/xlab/tablewriter"
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

func GetStatusTable() string {
	str := fmt.Sprintln("\nAverage LastBlock ", atomic.LoadInt64(&LastBlock.Dig), "\n")
	table := tablewriter.CreateTable()

	table.AddHeaders("Status", "Hostname", "LastBlock", "Weight", "LastUpdate")
	for _, server := range backends {
		table.AddRow(server.FSM.Current(), server.Hostname, server.LastBlock, server.Weight, time.Unix(server.TimeUpdate, 0).Format(time.RFC3339))
	}

	return str + table.Render()
}
