package scheduler

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type EthBlock struct {
	Dig   int64 `json:"block"`
	Hex   string
	Token string `json:"token"`
	Mutex sync.RWMutex
}

func GetTargetNode(backends map[string]Upstream, block int64, lastblock *EthBlock) (*url.URL, error) {

	//TODO: учитывать количество коннектов на бекенд
	//TODO: если количество коннектов исчерпано перейти на сервера бэкапа

	wBlock := block
	roulete := make([]string, 0)
	var srv Upstream
	if block == -1 {
		wBlock = atomic.LoadInt64(&lastblock.Dig)
	}

	for key, srv := range backends {
		if srv.FSM.Current() == "active" {
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
	if len(roulete) > 0 {
		rand.Seed(time.Now().Unix())
		winner := rand.Int() % len(roulete)
		srv = backends[roulete[winner]]
		log.Println("target node ", srv.Target, " for block:", fmt.Sprintf("%v", wBlock))
		return srv.GetURL()
	}

	return nil, errors.New("Not found avtive servers")

}
