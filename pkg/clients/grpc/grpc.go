package grpc

import (
	"context"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/server"

	"google.golang.org/grpc"
)

var _ clients.Client = (*Client)(nil)

type Client struct {
}

func (p *Client) Call(node *node.Node, in *message.Request) (*message.Response, error) {
	cc, err := grpc.Dial(node.Value)
	if err != nil {
		return nil, err
	}

	trellisClient := server.NewTrellisClient(cc)

	// TODO Context
	return trellisClient.Call(context.Background(), in)
}

func NewClient() (*Client, error) {
	return &Client{}, nil
}
