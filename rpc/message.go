package rpc

import (
	"encoding/json"
	"github.com/versus/gethinx/lib"
	"reflect"
	"errors"
)

type JsonRpcMessage struct {
	Version string          `json:"jsonrpc"`
	ID      int 		`json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  []json.RawMessage `json:"params,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (req JsonRpcMessage) GetStringParams(index int) (string,error) {
	var ret string
	var err error
	if index < len(req.Params) {
		if reflect.TypeOf(string(req.Params[index])).String() == "string" {
			ret = lib.TrimQuote(string(req.Params[index]))
		} else {
			err = errors.New("params is not string")
		}
	} else {
		err = errors.New("index out of range")
	}
	return ret, err
}

