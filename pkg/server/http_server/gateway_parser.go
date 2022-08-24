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

package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/errcode"
)

type GatewayParser struct {
	services map[string]*service.Service
}

type GatewayResponse struct {
	Code    int64  `json:"code"`
	TraceId string `json:"trace_id"`
	Msg     string `json:"msg,omitempty"`
	Payload []byte `json:"payload"`
}

func NewGatewayParser(services map[string]*service.Service) Parser {
	fmt.Println("gateway services", services)
	return &GatewayParser{
		services: services,
	}
}

func (p *GatewayParser) ParseRequest(ctx *fiber.Ctx) (*message.Request, error) {

	req := &message.Request{
		Payload: &message.Payload{
			Header: map[string]string{},
		},
	}

	urlPath := string(ctx.Request().URI().Path())
	if req.GetPayload().Header[mime.HeaderKeyRequestURIPath] == "" {
		req.GetPayload().Header[mime.HeaderKeyRequestURIPath] = urlPath
	}

	method := ctx.Method()
	if req.GetPayload().Header[mime.HeaderKeyRequestURIMethod] == "" {
		req.GetPayload().Header[mime.HeaderKeyRequestURIMethod] = method
	}

	if req.GetPayload().Header[mime.HeaderKeyUserAgent] == "" {
		req.GetPayload().Header[mime.HeaderKeyUserAgent] = string(ctx.Request().Header.UserAgent())
	}

	if req.GetPayload().Header[mime.HeaderKeyRequestURIQuery] == "" {
		req.GetPayload().Header[mime.HeaderKeyRequestURIQuery] = string(ctx.Request().URI().QueryString())
	}
	ct := string(ctx.Request().Header.Peek(mime.HeaderKeyContentType))
	req.GetPayload().Header[mime.HeaderKeyContentType] = ct

	clientIp := ClientIP(ctx.Context())
	if req.GetPayload().Header[mime.HeaderKeyClientIP] == "" {
		req.GetPayload().Header[mime.HeaderKeyClientIP] = clientIp
	}
	req.GetPayload().Header[mime.HeaderKeyRequestIP] = clientIp

	body := ctx.Request().Body()
	if body != nil {
		if err := req.GetPayload().SetBody(body); err != nil {
			return nil, err
		}
	}

	fmt.Println(urlPath, method)

	reqService, ok := p.services[queryFullpath(urlPath, method)]
	if !ok {
		return nil, errcode.Newf("not found handler: %s, %s", urlPath, method)
	}
	req.Service = reqService

	return req, nil
}

func (p *GatewayParser) ParseResponse(ctx *fiber.Ctx, req *message.Request, msg *message.Response) error {
	return parseGatewayResponse(ctx, req, msg)
}

func parseGatewayResponse(ctx *fiber.Ctx, req *message.Request, msg *message.Response) error {
	var ct string
	if msg != nil {
		ct = msg.GetPayload().Get(mime.HeaderKeyContentType)
	}
	if ct == "" {
		ct = mime.ContentTypeJson
	}

	c := codec.Select(ct)
	if c == nil {
		return errcode.Newf("not supported content type: %s", ct)
	}

	resp := &GatewayResponse{
		TraceId: req.GetPayload().Get(mime.HeaderKeyTraceID),
	}

	resp.Code = msg.Code
	if resp.Code != 0 {
		resp.Msg = "not ok"
	}

	resp.Payload = msg.GetPayload().GetBody()

	bs, _ := json.Marshal(resp)

	ctx.Set("Content-Type", ct)
	ctx.Response().SetStatusCode(http.StatusOK)
	ctx.Response().SetBody(bs)
	return nil
}
