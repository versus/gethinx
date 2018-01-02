package rpc

import "encoding/json"

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

