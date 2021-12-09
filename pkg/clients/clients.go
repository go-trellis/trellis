package clients

import (
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
)

type Client interface {
	Call(node *node.Node, in *message.Request) (*message.Response, error)
}
