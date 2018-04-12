package scheduler

import (
	"errors"
	"reflect"

	"github.com/versus/gethinx/lib"
)

// GetStringParams return index string from params field of request
func (req JSONRPCMessage) GetStringParams(index int) (string, error) {
	var ret string
	var err error
	if index < len(req.Params) {
		if reflect.ValueOf(req.Params[index]).IsNil() {
			err = errors.New("params is not string")
			return ret, err
		}
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
