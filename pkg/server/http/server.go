package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"trellis.tech/trellis.v1/pkg/clients/client"
	"trellis.tech/trellis.v1/pkg/clients/local"
	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/trellis"

	routing "github.com/go-trellis/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"trellis.tech/trellis/common.v0/errcode"
)

var (
	_ server.Server       = (*Server)(nil)
	_ component.Component = (*Server)(nil)
)

type Server struct {
	conf trellis.ServerConfig `yaml:"server_config" json:"server_config"`

	fastServer *fasthttp.Server
	fastRouter *routing.Router

	routes router.Router
}

func NewServer(conf trellis.ServerConfig) (*Server, error) {
	s := &Server{
		conf: conf,

		routes: router.NewRouter(conf.RouterConfig),

		fastRouter: routing.New(),
	}

	if err := s.routes.Start(); err != nil {
		return nil, err
	}

	return s, nil
}

type Handler struct {
	Method  string
	Path    string
	Uses    []routing.Handler
	Handler routing.Handler
}

func (p *Server) RegisterHandler(handlers ...Handler) {
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

func (p *Server) RegisterGroup(groupPath string, handlers ...Handler) {
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

	// TODO config to new component
	for _, comp := range p.conf.Components {
		if err := router.NewComponent(comp); err != nil {
			return err
		}
	}

	p.fastServer = &fasthttp.Server{
		Handler: p.fastRouter.HandleRequest,

		DisableKeepalive: true,
		//CloseOnShutdown:  true,
		IdleTimeout: time.Second * 30,
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
	if err := router.StopComponents(); err != nil {
		// TODO log
		fmt.Println(err)
	}

	if err := p.routes.Stop(); err != nil {
		// TODO log
		fmt.Println(err)
	}

	if err := p.fastServer.Shutdown(); err != nil {
		// TODO log
		fmt.Println(err)
		return err
	}

	return nil
}

func (p *Server) HandleHTTP(ctx *routing.Context) error {
	req, err := p.parseToRequest(ctx)
	if err != nil {
		return err
	}

	fmt.Println(*req)
	resp, err := p.Call(context.Background(), req)
	if err != nil {
		return err
	}

	return p.parseToResponse(ctx, resp)
}

func (p *Server) Call(ctx context.Context, msg *message.Request) (*message.Response, error) {
	// TODO with keys
	serviceNode, ok := p.routes.GetServiceNode(msg.GetService())
	if !ok {
		c, _ := local.NewClient()
		return c.Call(ctx, msg)
	}

	c, err := client.New(serviceNode)
	if err != nil {
		return nil, err
	}
	// TODO Options
	return c.Call(ctx, msg)
}

func (p *Server) Route(topic string, msg *message.Payload) (interface{}, error) {
	return nil, nil
}

func (*Server) parseToRequest(ctx *routing.Context) (*message.Request, error) {

	ct := string(ctx.Request.Header.Peek("Content-Type"))

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

	req.Payload = &message.Payload{
		Header: map[string]string{"Content-Type": ct},
	}
	req.Payload.Body = body

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

	ctx.SetContentType(codec.ContentTypeJson)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(bs)

	return nil
}
