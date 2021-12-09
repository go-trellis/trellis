package quic

import (
	"net/http"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
)

var _ clients.Client = (*Client)(nil)

type Client struct {
	client *http.Client
}

func (p *Client) Call(node *node.Node, in *message.Request) (*message.Response, error) {

	//p.client.Do()
	return nil, nil
}

func NewClient() (*Client, error) {
	var qconf quic.Config

	roundTripper := &http3.RoundTripper{
		QuicConfig: &qconf,
	}

	return &Client{
		client: &http.Client{
			Transport: roundTripper,
		},
	}, nil
}
