package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis.v1/pkg/node"

	"trellis.tech/trellis/common.v1/errcode"
)

var _ clients.Client = (*Client)(nil)

type Client struct {
	hc   *http.Client
	nd   *node.Node
	urlP *url.URL
}

func (p *Client) Call(ctx context.Context, in *message.Request, _ ...clients.CallOption) (*message.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	bs, _ := json.Marshal(in)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.urlP.String(), bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}

	req.Header.Set(mime.HeaderKeyContentType, mime.ContentTypeJson)

	hResp, err := p.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer hResp.Body.Close()

	body, err := ioutil.ReadAll(hResp.Body)
	if err != nil {
		return nil, err
	}

	if hResp.StatusCode != 200 {
		return nil, errcode.Newf("status code not 200, but %d, body: %s", hResp.StatusCode, string(body))
	}

	msg := &message.Response{}
	err = json.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// NewClient construct http instance
// TODO node.Metadata to http client config
func NewClient(nd *node.Node) (*Client, error) {
	if nd == nil {
		return nil, errcode.New("nil node")
	}
	urlP, err := url.Parse(getURL(nd.Value))
	if err != nil {
		return nil, err
	}
	return &Client{
		hc:   &http.Client{},
		nd:   nd,
		urlP: urlP,
	}, nil
}

func getURL(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "http://" + url
}
