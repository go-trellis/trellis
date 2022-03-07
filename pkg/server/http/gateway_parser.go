package http

import (
	"fmt"

	routing "github.com/go-trellis/fasthttp-routing"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis/common.v1/errcode"
)

type GatewayParser struct {
	services map[string]*service.Service
}

func NewGatewayParser(services map[string]*service.Service) Parser {
	return &GatewayParser{
		services: services,
	}
}

func (p *GatewayParser) ParseRequest(ctx *routing.Context) (*message.Request, error) {

	req := &message.Request{
		Payload: &message.Payload{
			Header: map[string]string{},
		},
	}

	urlPath := string(ctx.URI().Path())
	if req.GetPayload().Header[mime.HeaderKeyRequestURIPath] == "" {
		req.GetPayload().Header[mime.HeaderKeyRequestURIPath] = urlPath
	}

	method := string(ctx.Method())
	if req.GetPayload().Header[mime.HeaderKeyRequestURIMethod] == "" {
		req.GetPayload().Header[mime.HeaderKeyRequestURIMethod] = method
	}

	if req.GetPayload().Header[mime.HeaderKeyUserAgent] == "" {
		req.GetPayload().Header[mime.HeaderKeyUserAgent] = string(ctx.Request.Header.UserAgent())
	}

	if req.GetPayload().Header[mime.HeaderKeyRequestURIQuery] == "" {
		req.GetPayload().Header[mime.HeaderKeyRequestURIQuery] = string(ctx.URI().QueryString())
	}
	ct := string(ctx.Request.Header.Peek(mime.HeaderKeyContentType))
	req.GetPayload().Header[mime.HeaderKeyContentType] = ct

	clientIp := ClientIP(ctx.RequestCtx)
	if req.GetPayload().Header[mime.HeaderKeyClientIP] == "" {
		req.GetPayload().Header[mime.HeaderKeyClientIP] = clientIp
	}
	req.GetPayload().Header[mime.HeaderKeyRequestIP] = clientIp

	body := ctx.Request.Body()
	if body != nil {
		if err := req.GetPayload().SetBody(body); err != nil {
			return nil, err
		}
	}

	fmt.Println(string(body))

	reqService, ok := p.services[queryFullpath(urlPath, method)]
	if !ok {
		return nil, errcode.Newf("not found handler: %s, %s", urlPath, method)
	}
	req.Service = reqService

	return req, nil
}

func (p *GatewayParser) ParseResponse(ctx *routing.Context, msg *message.Response) error {
	return parseResponse(ctx, msg)
}
