package client

import (
	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/clients/grpc"
	"trellis.tech/trellis.v1/pkg/clients/http"
	"trellis.tech/trellis.v1/pkg/clients/quic"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis/common.v0/errcode"
)

// New TODO 解决client复用问题
func New(protocol node.Protocol) (clients.Client, error) {
	switch protocol {
	case node.Protocol_GRPC:
		return grpc.NewClient()
	case node.Protocol_HTTP:
		return http.NewClient()
	case node.Protocol_QUIC:
		return quic.NewClient()
	}

	return nil, errcode.Newf("not supported node protocol: %s", protocol)
}
