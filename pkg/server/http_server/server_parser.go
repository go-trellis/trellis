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
	"net/http"

	"github.com/gofiber/fiber/v2"

	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis/common.v1/errcode"
)

type ServerParser struct{}

func NewServerParser() Parser {
	return (*ServerParser)(nil)
}

func (*ServerParser) ParseRequest(ctx *fiber.Ctx) (*message.Request, error) {

	ct := string(ctx.Request().Header.ContentType())
	if ct == "" {
		ct = mime.ContentTypeJson
	}
	c := codec.Select(ct)
	if c == nil {
		return nil, errcode.Newf("not supported content type: %s", ct)
	}

	req := &message.Request{}

	body := ctx.Request().Body()
	if body != nil {
		if err := c.Unmarshal(ctx.Request().Body(), req); err != nil {
			return nil, err
		}
	}

	if req.GetPayload() == nil {
		req.Payload = &message.Payload{}
	}
	if req.GetPayload().GetHeader() == nil {
		req.GetPayload().Header = map[string]string{}
	}

	clientIp := ClientIP(ctx.Context())
	req.GetPayload().Header[mime.HeaderKeyContentType] = ct
	if req.GetPayload().Header[mime.HeaderKeyClientIP] == "" {
		req.GetPayload().Header[mime.HeaderKeyClientIP] = clientIp
	}
	req.GetPayload().Header[mime.HeaderKeyRequestIP] = clientIp

	return req, nil
}

func (*ServerParser) ParseResponse(ctx *fiber.Ctx, req *message.Request, msg *message.Response) error {
	return parseResponse(ctx, req, msg)
}

func parseResponse(ctx *fiber.Ctx, req *message.Request, msg *message.Response) error {
	ct := ""
	if msg == nil {
		msg = &message.Response{}
	} else {
		ct = msg.GetPayload().Get(mime.HeaderKeyContentType)
	}
	if ct == "" {
		ct = mime.ContentTypeJson
	}

	c := codec.Select(ct)
	if c == nil {
		return errcode.Newf("not supported content type: %s", ct)
	}

	bs, err := c.Marshal(msg)
	if err != nil {
		return err
	}

	msg.GetPayload().Set("Content-Type", ct)

	ctx.Set("Content-Type", ct)
	ctx.Response().SetStatusCode(http.StatusOK)
	ctx.Response().SetBody(bs)

	return nil
}
