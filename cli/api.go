package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/versus/gethinx"
	"github.com/versus/gethinx/backend"
	"github.com/versus/gethinx/ethblock"
	"github.com/versus/gethinx/middle"
	"github.com/versus/gethinx/proxyserver"
	"log"
	"net/http"
	"time"
)

func StartApi(addr string, addrAdmin string) {
	ar := gin.New()
	ar.LoadHTMLGlob("templates/*")

	ar.GET("/status", func(c *gin.Context) {
		c.HTML(http.StatusOK, "status.tmpl", gin.H{
			"title":     "Gethinx status page",
			"ws_server": conf.WebSocket,
		})
	})
	ar.GET("/ws", gethinx.WebSocketAdmin)
	ar.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	ar.POST("/api/v1/newblock", backend.SetBlock)
	ar.GET("/api/v1/status", gethinx.GetStatus)
	ar.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	go func() {
		err := ar.Run(addrAdmin)
		if err != nil {
			log.Println("Error run admin router: ", err.Error())
		}
	}()

	router := gin.Default()
	router.Use(middle.RequestLogger())
	router.Use(middle.ResponseLogger)

	router.Any("/", proxyserver.ReverseProxyHandler)

	err := router.Run(addr)
	if err != nil {
		log.Println("Error run gin router: ", err.Error())
	}

}

func AgentTickerUpstream() {
	log.Println(conf.Ticker)
	tick := time.Tick(time.Second * time.Duration(conf.Ticker))
	for {
		select {
		case <-tick:
			for key, srv := range *backends {
				if srv.FSM.Current() == "active" {
					lastTimeUpdate := time.Unix(srv.TimeUpdate, 0)
					now := time.Now()
					diff := now.Sub(lastTimeUpdate)
					//log.Println(srv.Target, "time sub is ", int64(diff/1000000000), "suspend time is ", int64(conf.Suspend))
					if int64(diff/1000000000) > int64(conf.Suspend) {
						srv.Mutex.Lock()
						srv.FSM.Event("suspend")
						srv.Mutex.Unlock()
					}
					bb := *backends
					bb[key] = srv
				}
			}
			checkAlive()

		}
	}
}

func checkAlive() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for key, srv := range *backends {
		srv.TargetLastBlock(ctx, &ethblock.LastBlock)
		//log.Println("checkAlive: ", srv.Target, " is ", srv.FSM.Current())
		bb := *backends
		bb[key] = srv
	}
	proxyserver.GenerateLastBlockAverage()
}
