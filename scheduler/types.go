package scheduler

import (
	"encoding/json"
	"sync"
)

type EthBlock struct {
	Dig        int64 `json:"block"`
	Hex        string
	Token      string `json:"token"`
	TimeUpdate int64  `json:"lastupdate"`
	Mutex      sync.RWMutex
}

// JSONRPCMessage structure for mapping request from ethclient
type JSONRPCMessage struct {
	Version string            `json:"jsonrpc"`
	ID      int               `json:"id,omitempty"`
	Method  string            `json:"method,omitempty"`
	Params  []json.RawMessage `json:"params,omitempty"`
	Error   *jsonError        `json:"error,omitempty"`
	Result  json.RawMessage   `json:"result,omitempty"`
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
