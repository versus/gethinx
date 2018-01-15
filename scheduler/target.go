package scheduler

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
)

type EthBlock struct {
	Dig   int64 `json:"block"`
	Hex   string
	Token string `json:"token"`
	Mutex sync.RWMutex
}

func GetTargetNode(backends map[string]Upstream, block int64, lastblock *EthBlock) (*url.URL, error) {
	wBlock := block
	if block == -1 {
		wBlock = atomic.LoadInt64(&lastblock.Dig)
	}
	log.Println("target node for block:", fmt.Sprintf("%v", wBlock))

	srv := backends["Q!@W#E$R%T^Y"]

	return srv.GetURL()
}
