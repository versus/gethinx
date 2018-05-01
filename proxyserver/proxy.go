package proxyserver

import (
	"encoding/json"
	"log"

	"net/url"

	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/backend"
	"github.com/versus/gethinx/buffer"
	"github.com/versus/gethinx/config"
	"github.com/versus/gethinx/ethblock"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Reversproxy struct {
	Backends backend.BackendList
	Conf     config.Config
	Mutex    sync.RWMutex
}

var instance *Reversproxy
var once sync.Once

func GetInstance() *Reversproxy {
	once.Do(func() {
		instance = &Reversproxy{}
	})
	return instance
}

func ReverseProxy(c *gin.Context) {

	var req ethblock.JSONRPCMessage
	//var block int64 = -1

	//var url *url.URL
	//var err error

	myreq := buffer.ReadRequestBody(c.Request.Body)
	c.Request.Body = myreq.Request

	bytes := []byte(myreq.Body)
	if err := json.Unmarshal(bytes, &req); err != nil {
		log.Println("Error unmarshal ", err.Error())
	}

	log.Println(req.Method)
	/*
		LastBlock.Mutex.RLock()
		hex := LastBlock.Hex
		LastBlock.Mutex.RUnlock()
	*/
	if req.Method == "eth_blockNumber" {
		c.JSON(200, gin.H{
			"jsonrpc": "2.0",
			"id":      req.ID,
			//"result":  hex,
		})
		return
	}

	if req.Method == "eth_getBlockByNumber" {
		// Req:  {"jsonrpc":"2.0","id":13,"method":"eth_getBlockByNumber","params":[null,false]}
		// Req:  {"jsonrpc":"2.0","id":14,"method":"eth_getBlockByNumber","params":["latest",false]}
		/*
		   		hexblock, err := req.GetStringParams(0)
		   		if err != nil {
		   			log.Println("Error get Params ", err.Error())
		   		}
		   /*
		   		/*
		   			if hexblock != "latest" {
		   				block, err = buffer.H2I(hexblock)
		   				if err != nil {
		   					log.Println("Error unhex block number ", err.Error())
		   				}
		   			}
		*/

	}
	/*
		url, err = GetTargetNode(backends, block, &LastBlock, conf.MaxResponseTime)
		if err != nil {
			log.Println("Error get URL for ReverseProxy  ", err.Error())
			c.Writer.WriteHeader(http.StatusBadGateway)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(c.Writer, c.Request)
	*/

}

func GetTargetNode(backends map[string]backend.Upstream, block int64, lastblock *ethblock.EthBlock, maxresponsetime int64) (*url.URL, error) {

	//TODO: учитывать количество коннектов на бекенд
	//TODO: если количество коннектов исчерпано перейти на сервера бэкапа

	wBlock := block
	roulete := make([]string, 0)
	var srv backend.Upstream
	if block == -1 {
		wBlock = atomic.LoadInt64(&lastblock.Dig)
	}

	for key, srv := range backends {
		if srv.FSM.Current() == "active" {
			if srv.ResponseTime < int64(time.Duration(maxresponsetime)*time.Millisecond) {
				if srv.LastBlock >= block {
					roulete = append(roulete, key)
					if srv.Weight > 1 {
						for i := 0; i < int(srv.Weight-1); i++ {
							roulete = append(roulete, key)
						}
					}
				}
			}

		}
	}
	if len(roulete) > 0 {
		rand.Seed(time.Now().Unix())
		winner := rand.Int() % len(roulete)
		srv = backends[roulete[winner]]
		log.Println("target node ", srv.Target, " for block:", fmt.Sprintf("%v", wBlock))
		return srv.URL()
	}

	return nil, errors.New("Not found avtive servers")

}
