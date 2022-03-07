package http

import (
	"net/http"

	routing "github.com/go-trellis/fasthttp-routing"
	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis/common.v1/errcode"
)

type ServerParser struct{}

func NewServerParser() Parser {
	return (*ServerParser)(nil)
}

func (*ServerParser) ParseRequest(ctx *routing.Context) (*message.Request, error) {

	ct := string(ctx.Request.Header.ContentType())
	if ct == "" {
		ct = mime.ContentTypeJson
	}
	c := codec.Select(ct)
	if c == nil {
		return nil, errcode.Newf("not supported content type: %s", ct)
	}

	req := &message.Request{}

	body := ctx.Request.Body()
	if body != nil {
		if err := c.Unmarshal(ctx.Request.Body(), req); err != nil {
			return nil, err
		}
	}

	if req.GetPayload() == nil {
		req.Payload = &message.Payload{}
	}
	if req.GetPayload().GetHeader() == nil {
		req.GetPayload().Header = map[string]string{}
	}

	clientIp := ClientIP(ctx.RequestCtx)
	req.GetPayload().Header[mime.HeaderKeyContentType] = ct
	if req.GetPayload().Header[mime.HeaderKeyClientIP] == "" {
		req.GetPayload().Header[mime.HeaderKeyClientIP] = clientIp
	}
	req.GetPayload().Header[mime.HeaderKeyRequestIP] = clientIp

	return req, nil
}

func (*ServerParser) ParseResponse(ctx *routing.Context, msg *message.Response) error {
	return parseResponse(ctx, msg)
}

func parseResponse(ctx *routing.Context, msg *message.Response) error {
	ct := string(ctx.Request.Header.Peek("Content-Type"))

	if msg == nil {
		msg = &message.Response{}
	} else {
		tmpCT := msg.GetPayload().Get(mime.HeaderKeyContentType)
		if tmpCT != "" {
			ct = tmpCT
		}
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
	ctx.SetContentType(ct)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(bs)

	return nil
}
