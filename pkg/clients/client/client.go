package client

import (
	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/clients/grpc"
	"trellis.tech/trellis.v1/pkg/clients/http"
	"trellis.tech/trellis.v1/pkg/clients/local"
	"trellis.tech/trellis.v1/pkg/clients/quic"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis/common.v0/errcode"
)

// New TODO 解决client复用问题
func New(n *node.Node, opts ...clients.NewOption) (clients.Client, error) {
	options := clients.NewOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	switch n.Protocol {
	case node.Protocol_LOCAL:
		return local.NewClient()
	case node.Protocol_GRPC:
		return grpc.NewClient(n)
	case node.Protocol_HTTP:
		return http.NewClient(n)
	case node.Protocol_QUIC:
		return quic.NewClient(n)
	}

	return nil, errcode.Newf("not supported node protocol: %s", n.GetProtocol())
}
