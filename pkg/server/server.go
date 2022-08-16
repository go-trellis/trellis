/*
Copyright © 2022 Henry Huang <hhh@rutcode.com>
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package server

import (
	"context"

	routing "github.com/go-trellis/fasthttp-routing"
	"google.golang.org/grpc"
	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis/common.v1/errcode"
)

type Server interface {
	lifecycle.LifeCycle
}

type Caller interface {
	// Call allows a single request to be made
	Call(ctx context.Context, in *message.Request, opts ...CallOption) (*message.Response, error)
	//// Publish publishes a payload and returns an empty payload
	//Publish(ctx context.Context, in *message.Request, opts ...CallOption) (*message.Response, error)
}

type CallOption func(*CallOptions)

type CallOptions struct {
	GRPCCallOptions []grpc.CallOption
}

func GRPCCallOption(opts ...grpc.CallOption) CallOption {
	return func(options *CallOptions) {
		options.GRPCCallOptions = append(options.GRPCCallOptions, opts...)
	}
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
