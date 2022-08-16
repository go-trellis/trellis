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

package local

import (
	"context"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis.v1/pkg/server"

	"trellis.tech/trellis/common.v1/errcode"
	"trellis.tech/trellis/common.v1/json"
)

var (
	_ server.Caller = (*Client)(nil)

	c = &Client{}
)

type Client struct{}

func NewClient() (server.Caller, []server.CallOption, error) {
	return c, nil, nil
}

func (*Client) Call(_ context.Context, in *message.Request, _ ...server.CallOption) (*message.Response, error) {
	comp := component.GetComponent(in.GetService())
	if comp == nil {
		return nil, errcode.Newf("not found component: %s", in.GetService().FullPath())
	}
	hResp, err := comp.Route(in.GetService().GetTopic(), in.GetPayload())
	if err != nil {
		// TODO log err
		return nil, err
	}

	if hResp == nil {
		return &message.Response{
			Code:    0,
			Payload: &message.Payload{Header: in.GetPayload().GetHeader()},
		}, nil
	}

	switch t := hResp.(type) {
	case message.Response:
		return &t, nil
	case *message.Response:
		return t, nil
	case *message.Payload:
		return &message.Response{
			Code:    0,
			Payload: t,
		}, nil
	case message.Payload:
		return &message.Response{
			Code:    0,
			Payload: &t,
		}, nil
	default:
		bs, err := json.Marshal(hResp)
		if err != nil {
			return nil, err
		}
		resp := &message.Response{
			Code: 0,
			Payload: &message.Payload{
				Header: in.GetPayload().GetHeader(),
				Body:   bs,
			},
		}

		resp.Payload.Header[mime.HeaderKeyContentType] = mime.ContentTypeJson
		return resp, nil
	}
}
