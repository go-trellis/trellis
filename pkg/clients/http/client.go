package http

import (
	"net/http"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
)

var _ clients.Client = (*Client)(nil)

type Client struct {
	client *http.Client
}

func (p *Client) Call(node *node.Node, in *message.Request) (*message.Response, error) {
	return nil, nil
}

func NewClient() (*Client, error) {

	return &Client{
		client: &http.Client{},
	}, nil
}
