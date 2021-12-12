package http

import (
	"context"
	"net/http"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
)

var _ clients.Client = (*Client)(nil)

type Client struct {
	client *http.Client
}

func (p *Client) Call(ctx context.Context, in *message.Request, opts ...clients.CallOption) (*message.Response, error) {
	return nil, nil
}

func NewClient(node *node.Node) (*Client, error) {
	return &Client{
		client: &http.Client{},
	}, nil
}
