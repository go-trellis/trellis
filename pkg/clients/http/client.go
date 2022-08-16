/*
Copyright © 2022 Henry Huang <hhh@rutcode.com>
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/server"

	"trellis.tech/trellis/common.v1/errcode"
)

var _ server.Caller = (*Client)(nil)

type Client struct {
	hc   *http.Client
	nd   *node.Node
	urlP *url.URL
}

func (p *Client) Call(ctx context.Context, in *message.Request, _ ...server.CallOption) (*message.Response, error) {
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
	defer func() {
		_, _ = io.Copy(io.Discard, hResp.Body)
		_ = hResp.Body.Close()
	}()

	body, err := io.ReadAll(hResp.Body)
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

// NewClient construct http_server instance
// TODO node.Metadata to http_server client config
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
