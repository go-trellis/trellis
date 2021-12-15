package clients

import (
	"context"

	"trellis.tech/trellis.v1/pkg/message"

	"google.golang.org/grpc"
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
