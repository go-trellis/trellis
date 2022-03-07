package http

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis.v1/pkg/tracing"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/trellis"

	"github.com/dgrr/http2"
	routing "github.com/go-trellis/fasthttp-routing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
)

var (
	_ server.Server = (*Server)(nil)
)

type Server struct {
	name string // Server Name

	conf *trellis.HTTPServerConfig

	fastServer *fasthttp.Server
	fastRouter *routing.Router

	router router.Router

	parser Parser

	tracing bool

	services map[string]*service.Service
}

type Parser interface {
	ParseRequest(ctx *routing.Context) (*message.Request, error)
	ParseResponse(ctx *routing.Context, msg *message.Response) error
}

func NewServer(opts ...Option) (*Server, error) {
	s := &Server{
		services: map[string]*service.Service{},

		fastRouter: routing.New(),
	}

	for _, o := range opts {
		o(s)
	}

	return s, nil
}

func (p *Server) Start() error {

	if p.router != nil {
		if err := p.router.Start(); err != nil {
			return err
		}
	}

	if p.fastRouter == nil {
		p.fastRouter = routing.New()
	}

	if err := p.registerGroupHandlers(p.fastRouter, p.conf.Groups...); err != nil {
		return err
	}

	if err := p.registerHandlers(p.fastRouter, p.conf.Handlers...); err != nil {
		return err
	}

	fmt.Println(p.services)

	p.fastServer = &fasthttp.Server{
		Handler:          p.fastRouter.HandleRequest,
		DisableKeepalive: p.conf.DisableKeepAlive,
		IdleTimeout:      p.conf.IdleTimeout,
		CloseOnShutdown:  true,
	}

	if p.conf.IsGateway {
		p.parser = NewGatewayParser(p.services)
	} else {
		p.parser = NewServerParser()
	}

	go func() {
		var err error

		switch p.conf.Protocol {
		case "http2":
			serverConfig := http2.ServerConfig{
				PingInterval:         p.conf.HTTP2Config.PingInterval,
				MaxConcurrentStreams: p.conf.HTTP2Config.MaxConcurrentStreams,
				Debug:                p.conf.HTTP2Config.Debug,
			}
			http2.ConfigureServer(p.fastServer, serverConfig)
		}

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
	fmt.Println("in........")
	var span opentracing.Span

	if opentracing.IsGlobalTracerRegistered() {

		span, _ = opentracing.StartSpanFromContext(ctx, p.name)

		ext.HTTPUrl.Set(span, ctx.Request.URI().String())
		ext.HTTPMethod.Set(span, string(ctx.Request.Header.Method()))
		ext.Component.Set(span, p.name)

		defer func() {

			headers := http.Header{}
			ctx.RequestCtx.Request.Header.VisitAllInOrder(p.visitor(headers))
			opentracing.GlobalTracer().Inject(span.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(headers),
			)

			ext.HTTPStatusCode.Set(span, uint16(ctx.Response.StatusCode()))
			if err != nil {
				ext.Error.Set(span, true)
				span.LogFields(opentracingLog.String("handler call error", err.Error()))
			}
			span.Finish()
		}()
	}

	fmt.Println("parse type", reflect.TypeOf(p.parser))

	var req *message.Request
	if req, err = p.parser.ParseRequest(ctx); err != nil {
		return routing.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fmt.Println(*req)

	if req.GetPayload().GetHeader() == nil {
		req.GetPayload().Header = map[string]string{}
		req.GetPayload().Header[mime.HeaderKeyTraceID] = tracing.GetTraceID(span)
	}

	var resp *message.Response
	if resp, err = p.router.Call(ctx, req); err != nil {
		fmt.Println("call server failed", err)
		return
	}

	err = p.parser.ParseResponse(ctx, resp)

	return
}

func (*Server) visitor(header http.Header) func(key, value []byte) {
	return func(key, value []byte) {
		header.Add(string(key), string(value))
	}
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
