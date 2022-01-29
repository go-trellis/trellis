package quic

//import (
//	"context"
//	"net/http"
//
//	"trellis.tech/trellis.v1/pkg/clients"
//	"trellis.tech/trellis.v1/pkg/message"
//	"trellis.tech/trellis.v1/pkg/node"
//
//	"github.com/lucas-clemente/quic-go"
//	"github.com/lucas-clemente/quic-go/http3"
//)
//
//var _ clients.Client = (*Client)(nil)
//
//type Client struct {
//	client *http.Client
//}
//
//func (p *Client) Call(ctx context.Context, in *message.Request, option ...clients.CallOption) (*message.Response, error) {
//
//	//p.client.Do()
//	return nil, nil
//}
//
//func NewClient(n *node.Node) (*Client, error) {
//	var qconf quic.Config
//
//	roundTripper := &http3.RoundTripper{
//		QuicConfig: &qconf,
//	}
//
//	return &Client{
//		client: &http.Client{
//			Transport: roundTripper,
//		},
//	}, nil
//}
