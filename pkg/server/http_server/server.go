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
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"

	"trellis.tech/trellis.v1/pkg/clients/client"
	"trellis.tech/trellis.v1/pkg/clients/local"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/tracing"
	"trellis.tech/trellis.v1/pkg/trellis"

	"github.com/dgrr/http2"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/valyala/fasthttp"
	"trellis.tech/trellis/common.v1/logger"
)

var (
	_ server.Server = (*Server)(nil)
)

type Server struct {
	name string // Server Name

	conf *trellis.HTTPServerConfig

	fiberApp    *fiber.App
	fiberRouter fiber.Router

	router router.Router

	parser Parser

	tracing bool

	services map[string]*service.Service

	logger logger.Logger
}

type Parser interface {
	ParseRequest(ctx *fiber.Ctx) (*message.Request, error)
	ParseResponse(ctx *fiber.Ctx, req *message.Request, msg *message.Response) error
}

func NewServer(opts ...Option) (*Server, error) {
	s := &Server{
		services: map[string]*service.Service{},

		fiberApp: fiber.New(),
	}

	for _, o := range opts {
		o(s)
	}

	s.fiberRouter = s.fiberApp.Use(recover.New(recover.Config{EnableStackTrace: s.conf.RecoverTrace}))

	fmt.Println("p.fiberRouter", s.fiberRouter)

	return s, nil
}

func (p *Server) Start() error {

	if p.router != nil {
		if err := p.router.Start(); err != nil {
			return err
		}
	}

	if err := p.registerGroupHandlers(p.conf.Groups...); err != nil {
		return err
	}

	if err := p.registerHandlers(p.conf.Handlers...); err != nil {
		return err
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
			http2.ConfigureServer(p.fiberApp.Server(), serverConfig)
		}

		if p.conf.EnableTLS {
			err = p.fiberApp.ListenTLS(p.conf.Address, p.conf.TLSConfig.CertPath, p.conf.TLSConfig.KeyPath)
		} else {
			err = p.fiberApp.Listen(p.conf.Address)
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
		fmt.Println(err)
	}

	if err := p.router.Stop(); err != nil {
		// TODO log
		fmt.Println(err)
	}

	if err := p.fiberApp.Shutdown(); err != nil {
		// TODO log
		return err
	}

	return nil
}

func (p *Server) HandleHTTP(ctx *fiber.Ctx) (err error) {
	fmt.Println("in........")
	var span opentracing.Span

	httputil.DumpRequest(ctx.Request().SetURI(""))

	if opentracing.IsGlobalTracerRegistered() {

		span, _ = opentracing.StartSpanFromContext(ctx.Context(), p.name)

		ext.HTTPUrl.Set(span, ctx.Request().URI().String())
		ext.HTTPMethod.Set(span, string(ctx.Request().Header.Method()))
		ext.Component.Set(span, p.name)

		defer func() {
			headers := http.Header{}
			ctx.Request().Header.VisitAllInOrder(p.visitor(headers))
			if dErr := opentracing.GlobalTracer().Inject(span.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(headers),
			); dErr != nil {
				span.LogFields(opentracingLog.String("handler inject error", dErr.Error()))
			}

			ext.HTTPStatusCode.Set(span, uint16(ctx.Response().StatusCode()))
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
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		ctx.Response().SetBody([]byte(err.Error()))
		return
	}

	fmt.Println("*req：", req)

	if req.GetPayload().GetHeader() == nil {
		req.GetPayload().Header = map[string]string{}
	}
	req.GetPayload().Header[mime.HeaderKeyTraceID] = tracing.GetTraceID(span)

	var (
		c    server.Caller
		opts []server.CallOption
	)
	serviceNode, ok := p.router.GetServiceNode(req.GetService(), req.String())
	if !ok {
		c, opts, err = local.NewClient()
	} else {
		c, opts, err = client.New(serviceNode)
	}
	if err != nil {
		return err
	}

	fmt.Println("req...", req)
	resp, err := c.Call(ctx.Context(), req, opts...)
	if err != nil {
		return err
	}

	err = p.parser.ParseResponse(ctx, req, resp)

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
