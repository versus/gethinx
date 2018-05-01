package gethinx

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xlab/tablewriter"
)

func JsonBackends() []string {
	var backendsServerList []string
	/*
		for _, server := range backends {
			jsrv, err := json.Marshal(server)
			if err != nil {
				log.Print("Error marchal backendsServerList for status ", err.Error())
			}
			backendsServerList = append(backendsServerList, string(jsrv))
		}
	*/
	return backendsServerList
}

func GetStatus(c *gin.Context) {

	backendsServerList := JsonBackends()
	c.JSON(200, gin.H{
		"ServerList": backendsServerList,
		//"LastBlockAverage": atomic.LoadInt64(&LastBlock.Dig),
	})
}

func GetStatusTable() string {

	//str := fmt.Sprintln("\nAverage LastBlock ", atomic.LoadInt64(&LastBlock.Dig), "\n")
	str := fmt.Sprintln("\nAverage LastBlock  \n")
	table := tablewriter.CreateTable()

	table.AddHeaders("Status", "Hostname", "LastBlock", "Weight", "LastUpdate", "ResponseTime")

	//for _, server := range backends {
	//	table.AddRow(server.FSM.Current(), server.Hostname, server.LastBlock, server.Weight, time.Unix(server.TimeUpdate, 0).Format(time.RFC3339), server.ResponseTime)
	//}

	return str + table.Render()
}
