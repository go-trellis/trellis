package server

import (
	routing "github.com/go-trellis/fasthttp-routing"
	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis/common.v1/errcode"
)

type Server interface {
	lifecycle.LifeCycle
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
