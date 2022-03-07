package server

import (
	"context"

	routing "github.com/go-trellis/fasthttp-routing"
	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis/common.v1/errcode"
)

type Server interface {
	lifecycle.LifeCycle
}

type Caller interface {
	// Call allows a single request to be made
	Call(context.Context, *message.Request) (*message.Response, error)
	// Publish publishes a payload and returns an empty payload
	Publish(context.Context, *message.Request) error
}

var uses = make(map[string]routing.Handler)

func RegisterUseFunc(name string, rh routing.Handler) {
	_, ok := uses[name]
	if ok {
		panic(errcode.Newf("use function already exist: %s", name))
	}
	uses[name] = rh
}

func GetUseFunc(name string) (routing.Handler, error) {
	use, ok := uses[name]
	if !ok {
		panic(errcode.Newf("not found use function: %s", name))
	}

	return use, nil
}
