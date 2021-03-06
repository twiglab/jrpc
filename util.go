package jrpc

import (
	"context"
	"encoding/json"
)

const echoName = "echo"

type (
	Echo struct {
	}

	EchoParams struct {
		Name string `json:"name"`
	}

	EchoResult struct {
		Message string `json:"message"`
	}
)

func NewEcho() *Echo {
	return new(Echo)
}

func (e *Echo) Params() interface{} {
	return new(EchoParams)
}

func (e *Echo) Result() interface{} {
	return new(EchoResult)
}

func (e *Echo) Name() string {
	return echoName
}

func (h *Echo) ServeJSONRPC(c context.Context, params *json.RawMessage) (interface{}, *Error) {

	var p EchoParams
	if err := Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return &EchoResult{
		Message: "Hello, " + p.Name,
	}, nil
}
