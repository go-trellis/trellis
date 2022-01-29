package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"

	"trellis.tech/trellis/common.v1/errcode"
)

var _ clients.Client = (*Client)(nil)

type Client struct {
	hc *http.Client
	nd *node.Node
}

func (p *Client) Call(ctx context.Context, in *message.Request) (*message.Response, error) {

	if ctx == nil {
		ctx = context.Background()
	}
	bs, _ := json.Marshal(in)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, in.GetService().String(), bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}

	hResp, err := p.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer hResp.Body.Close()

	if hResp.StatusCode != 200 {
		return nil, errcode.Newf("status code not 200, but %d", hResp.StatusCode)
	}

	body, err := ioutil.ReadAll(hResp.Body)
	if err != nil {
		return nil, err
	}

	msg := &message.Response{}
	err = json.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// NewClient construct http instance
// TODO node.Metadata to http config
func NewClient(nd *node.Node) (*Client, error) {
	if nd == nil {
		return nil, errcode.New("nil node")
	}
	return &Client{
		hc: &http.Client{},
		nd: nd,
	}, nil
}
