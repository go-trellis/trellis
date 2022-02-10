package http

import (
	"log"
	"net/http"
	"strings"

	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/trellis"

	routing "github.com/go-trellis/fasthttp-routing"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
	"trellis.tech/trellis/common.v1/errcode"
)

var (
	_ server.Server       = (*Server)(nil)
	_ component.Component = (*Server)(nil)
)

type Server struct {
	name string // Server Name

	conf *trellis.HTTPServerConfig

	fastServer *fasthttp.Server
	fastRouter *routing.Router

	router router.Router

	tracing bool
}

func NewServer(opts ...Option) (*Server, error) {
	s := &Server{
		fastRouter: routing.New(),
	}

	for _, o := range opts {
		o(s)
	}

	var hs []*Handler
	for _, hCfg := range s.conf.Handlers {
		h, err := getHTTPHandler(hCfg)
		if err != nil {
			return nil, err
		}
		hs = append(hs, h)
	}
	s.RegisterHandler(hs...)

	for _, group := range s.conf.Groups {
		hgs, err := getHTTPGroupHandlers(group)
		if err != nil {
			return nil, err
		}
		s.RegisterGroup(group.Path, hgs...)
	}

	if s.router != nil {
		if err := s.router.Start(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

type Handler struct {
	Method  string
	Path    string
	Uses    []routing.Handler
	Handler routing.Handler
}

func (p *Server) RegisterHandler(handlers ...*Handler) {
	for _, handler := range handlers {
		var fastHandlers = handler.Uses
		if handler.Handler == nil {
			fastHandlers = append(fastHandlers, p.HandleHTTP)
		} else {
			fastHandlers = append(fastHandlers, handler.Handler)
		}
		p.fastRouter.To(handler.Method, handler.Path, fastHandlers...)
	}
}

func (p *Server) RegisterGroup(groupPath string, handlers ...*Handler) {
	group := p.fastRouter.Group(groupPath)
	for _, handler := range handlers {
		var fastHandlers = handler.Uses
		if handler.Handler == nil {
			fastHandlers = append(fastHandlers, p.HandleHTTP)
		} else {
			fastHandlers = append(fastHandlers, handler.Handler)
		}
		group.To(handler.Method, handler.Path, fastHandlers...)
	}
}

func (p *Server) Start() error {

	p.fastServer = &fasthttp.Server{
		Handler: p.fastRouter.HandleRequest,

		DisableKeepalive: p.conf.DisableKeepAlive,
		//CloseOnShutdown:  true,
		IdleTimeout: p.conf.IdleTimeout,
	}

	go func() {
		var err error
		if p.conf.EnableTLS {
			err = p.fastServer.ListenAndServeTLS(p.conf.Address, p.conf.TLSConfig.CertPath, p.conf.TLSConfig.KeyPath)
		} else {
			err = p.fastServer.ListenAndServe(p.conf.Address)
		}

		if err != nil {
			if err != http.ErrServerClosed {
				log.Fatalln(err)
			}
		}
	}()

	return nil
}

func (p *Server) Stop() error {
	if err := component.StopComponents(); err != nil {
		// TODO log
	}

	if err := p.router.Stop(); err != nil {
		// TODO log
	}

	if err := p.fastServer.Shutdown(); err != nil {
		// TODO log
		return err
	}

	return nil
}

func (p *Server) HandleHTTP(ctx *routing.Context) (err error) {

	var span opentracing.Span

	//opentracing.GlobalTracer()
	if opentracing.IsGlobalTracerRegistered() {

		span, _ = opentracing.StartSpanFromContext(ctx, p.name)

		ext.HTTPUrl.Set(span, ctx.Request.URI().String())
		ext.HTTPMethod.Set(span, string(ctx.Request.Header.Method()))
		ext.Component.Set(span, p.name)

		headers := http.Header{}
		ctx.RequestCtx.Request.Header.VisitAllInOrder(p.visitor(headers))
		opentracing.GlobalTracer().Inject(span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(headers),
		)

		defer func() {
			ext.HTTPStatusCode.Set(span, uint16(ctx.Response.StatusCode()))
			if err != nil {
				ext.Error.Set(span, true)
				span.LogFields(opentracingLog.String("handler call error", err.Error()))
			}
			span.Finish()
		}()
	}

	var msg *message.Request
	msg, err = p.parseToRequest(ctx, span)
	if err != nil {
		return err
	}

	var resp *message.Response
	resp, err = p.router.Call(ctx, msg)
	if err != nil {
		return err
	}

	err = p.parseToResponse(ctx, resp)
	return
}

func (*Server) visitor(header http.Header) func(key, value []byte) {
	return func(key, value []byte) {
		header.Add(string(key), string(value))
	}
}

func (p *Server) Route(_topic string, _payload *message.Payload) (interface{}, error) {
	return nil, nil
}

func (*Server) parseToRequest(ctx *routing.Context, span opentracing.Span) (*message.Request, error) {

	ct := string(ctx.Request.Header.Peek(mime.HeaderKeyContentType))

	clientIp := ClientIP(ctx.RequestCtx)

	c := codec.Select(ct)
	if c == nil {
		return nil, errcode.Newf("not supported content type: %s", ct)
	}

	req := &message.Request{}
	body := ctx.Request.Body()
	if body == nil {
		return req, nil
	}

	if err := c.Unmarshal(body, req); err != nil {
		return nil, err
	}

	if req.GetPayload() == nil {
		req.Payload = &message.Payload{}
	}
	if req.GetPayload().GetHeader() == nil {
		req.GetPayload().Header = map[string]string{}
	}

	if req.GetPayload().GetHeader()[mime.HeaderKeyTraceID] == "" {
		req.GetPayload().Header[mime.HeaderKeyTraceID] = uuid.NewString()
	}

	req.GetPayload().Header[mime.HeaderKeyContentType] = ct
	if req.GetPayload().Header[mime.HeaderKeyClientIP] == "" {
		req.GetPayload().Header[mime.HeaderKeyClientIP] = clientIp
	}
	req.GetPayload().Header[mime.HeaderKeyRequestIP] = clientIp
	req.GetPayload().Header[mime.HeaderKeyRequestID] = uuid.NewString()

	if span != nil {
		// todo added tags
	}
	return req, nil
}

func (*Server) parseToResponse(ctx *routing.Context, msg *message.Response) error {

	ct := string(ctx.Request.Header.Peek("Content-Type"))
	c := codec.Select(ct)
	if c == nil {
		return errcode.Newf("not supported content type: %s", ct)
	}

	if msg == nil {
		msg = &message.Response{}
	}
	msg.GetPayload().Set("Content-Type", ct)

	bs, err := c.Marshal(msg)
	if err != nil {
		return err
	}

	ctx.SetContentType(mime.ContentTypeJson)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(bs)

	return nil
}

// ClientIP 获取真实的IP  1.1.1.1, 2.2.2.2, 3.3.3.3
func ClientIP(ctx *fasthttp.RequestCtx) string {
	clientIP := string(ctx.Request.Header.Peek("X-Forwarded-For"))
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
		//获取最开始的一个 即 1.1.1.1
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(string(ctx.Request.Header.Peek("X-Real-Ip")))
	if len(clientIP) > 0 {
		return clientIP
	}
	return ctx.RemoteIP().String()
}
