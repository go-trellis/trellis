package client

import (
	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/clients/grpc"
	"trellis.tech/trellis.v1/pkg/clients/http"
	"trellis.tech/trellis.v1/pkg/clients/local"
	"trellis.tech/trellis.v1/pkg/node"

	"trellis.tech/trellis/common.v1/errcode"
)

// New TODO 解决client复用问题
func New(nd *node.Node) (clients.Client, error) {
	if nd == nil {
		return local.NewClient()
	}
	switch nd.GetProtocol() {
	case node.Protocol_LOCAL:
		return local.NewClient()
	case node.Protocol_GRPC:
		return grpc.NewClient(nd)
	case node.Protocol_HTTP:
		return http.NewClient(nd)
		//case node.Protocol_QUIC:
		//	return quic.NewClient(n)
	}

	return nil, errcode.Newf("not supported node protocol: %d, %s",
		nd.GetProtocol(), nd.GetProtocol().String())
}
