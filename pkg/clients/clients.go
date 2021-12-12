package clients

import (
	"context"

	"google.golang.org/grpc"
	"trellis.tech/trellis.v1/pkg/message"
)

type NewOption func(*NewOptions)
type NewOptions struct {
	GrpcOptions []grpc.DialOption
}

type CallOption func(*CallOptions)
type CallOptions struct {
	GrpcCallOptions []grpc.CallOption
}

type Client interface {
	Call(context.Context, *message.Request, ...CallOption) (*message.Response, error)
}
