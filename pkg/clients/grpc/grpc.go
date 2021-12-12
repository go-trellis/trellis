package grpc

import (
	"context"

	"google.golang.org/grpc"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/server"
)

var _ clients.Client = (*Client)(nil)

type Client struct {
	trellisClient server.TrellisClient
}

func (p *Client) Call(ctx context.Context, in *message.Request, opts ...clients.CallOption) (*message.Response, error) {
	options := &clients.CallOptions{}
	for _, o := range opts {
		o(options)
	}

	// TODO Context
	return p.trellisClient.Call(ctx, in, options.GrpcCallOptions...)
}

func NewClient(node *node.Node) (c *Client, err error) {
	client := &Client{}
	cc, err := grpc.Dial(node.Value, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client.trellisClient = server.NewTrellisClient(cc)
	return client, nil
}
