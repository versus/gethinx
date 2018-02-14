package main

import (
	"fmt"
	"net/http"

	"sync/atomic"

	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSStatusResponse struct {
	LastBlockAverage           int64
	LastBlockAverageTimeUpdate int64    `json:"lastupdate"`
	Upstream                   []string `json:"upstreams"`
}

func webSocketAdmin(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}
	answerStruct := &WSStatusResponse{
		LastBlockAverage: 0,
		Upstream:         []string{},
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error read websocket ", err.Error())
			break
		}
		if string(msg) == "status" {
			answerStruct.LastBlockAverage = atomic.LoadInt64(&LastBlock.Dig)
			answerStruct.Upstream = JsonBackends()
			answer, err := json.Marshal(answerStruct)
			if err != nil {
				log.Println("Error marshal answer ", err.Error())
			}
			conn.WriteMessage(t, answer)
		}
	}
}
