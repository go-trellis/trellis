package clients

import "google.golang.org/grpc"

type CallOption func(*CallOptions)
type CallOptions struct {
	GrpcCallOptions []grpc.CallOption
}

func GrpcCallOptions(opts []grpc.CallOption) CallOption {
	return func(options *CallOptions) {
		options.GrpcCallOptions = opts
	}
}
