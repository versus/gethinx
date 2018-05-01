package ethblock

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

type EthBlock struct {
	Dig        int64 `json:"block"`
	Hex        string
	Token      string `json:"token"`
	TimeUpdate int64  `json:"lastupdate"`
	Mutex      sync.RWMutex
}

var LastBlock EthBlock

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

// GetStringParams return index string from params field of request
func (req JSONRPCMessage) GetStringParams(index int) (string, error) {
	var ret string
	var err error
	if index < len(req.Params) {
		if reflect.TypeOf(string(req.Params[index])).String() == "string" {
			ret = trimQuote(string(req.Params[index]))
		} else {
			err = errors.New("params is not string")
		}
	} else {
		err = errors.New("index out of range")
	}
	return ret, err
}

func trimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

/*
import (
"sync"
)

type singleton struct {
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

*/
