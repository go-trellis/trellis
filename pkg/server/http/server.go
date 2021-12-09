package http

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"trellis.tech/trellis.v1/pkg/clients/client"
	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/service"

	routing "github.com/go-trellis/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"trellis.tech/trellis/common.v0/errcode"
)

var (
	_ server.Server       = (*Server)(nil)
	_ component.Component = (*Server)(nil)
)

type Server struct {
	fastServer *fasthttp.Server
	fastRouter *routing.Router

	conf Config

	routes router.Router

	compManager component.Manager
}

type Config struct {
	Address   string `yaml:"address" json:"address"`
	EnableTLS bool   `yaml:"enable_tls" json:"enable_tls"`
	CertFile  string `yaml:"cert_file" json:"cert_file"`
	KeyFile   string `yaml:"key_file" json:"key_file"`

	RouterConfig router.Config `json:"router_config" yaml:"router_config"`

	Services []*service.Service `yaml:"services" json:"services"`
}

func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	cfg.ParseFlagsWithPrefix(f, "")
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *Config) ParseFlagsWithPrefix(f *flag.FlagSet, prefix string) {
	f.StringVar(&cfg.Address, prefix+".server.address", "", "")
	f.BoolVar(&cfg.EnableTLS, prefix+".server.enable_tls", false, "")
	f.StringVar(&cfg.CertFile, prefix+".server.cert_file", "", "")
	f.StringVar(&cfg.KeyFile, prefix+".server.cert_file", "", "")
	cfg.RouterConfig.ParseFlagsWithPrefix(f, prefix)
}

func NewServer(conf Config) (*Server, error) {
	s := &Server{
		conf: conf,

		fastRouter:  routing.New(),
		compManager: component.GetManager(),
	}

	s.routes = router.NewRouter(s.conf.RouterConfig)
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
	for _, s := range p.conf.Services {
		if err := p.compManager.NewComponent(s); err != nil {
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
			err = p.fastServer.ListenAndServeTLS(p.conf.Address, p.conf.CertFile, p.conf.KeyFile)
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
	if err := p.compManager.Stop(); err != nil {
		// TODO log
		fmt.Println(err)
	}
	//p.compManager.
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

func (p *Server) Register(s *service.ServiceNode) error {
	return p.routes.Register(s)
}

func (p *Server) Deregister(s *service.ServiceNode) error {
	return p.routes.Deregister(s)
}

func (p *Server) Watch(s *service.Service) error {
	return p.routes.Watch(s)
}

func (p *Server) HandleHTTP(ctx *routing.Context) error {
	req, err := p.parseToRequest(ctx)
	if err != nil {
		return err
	}

	fmt.Println(*req)
	resp, err := p.Handle(context.Background(), req)
	if err != nil {
		return err
	}

	return p.parseToResponse(ctx, resp)
}

func (p *Server) Handle(ctx context.Context, msg *message.Request) (*message.Response, error) {
	comp := p.compManager.GetComponent(msg.GetService())
	if comp != nil {
		hResp, err := comp.Route(msg.GetPayload())
		if err != nil {
			// TODO log err
			return nil, err
		}

		if hResp == nil {
			return &message.Response{
				Code: 0,
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
			return &message.Response{
				Code: 0,
				Payload: &message.Payload{
					Header: map[string]string{"Content-Type": codec.ContentTypeJson},
					Body:   bs,
				},
			}, nil
		}
	}

	// TODO with keys
	node, ok := p.routes.GetServiceNode(msg.GetService())
	if !ok {
		return nil, errcode.Newf("not found handler: %+v", msg.GetService())
	}

	c, err := client.New(node.Protocol)
	if err != nil {
		return nil, err
	}

	return c.Call(node, msg)
}

func (p *Server) Route(msg *message.Payload) (interface{}, error) {
	return nil, nil
}

func (p *Server) parseToRequest(ctx *routing.Context) (*message.Request, error) {

	ct := string(ctx.Request.Header.Peek("Content-Type"))

	c := codec.Select(ct)
	if c == nil {
		return nil, errcode.Newf("not supported content type: %s", ct)
	}

	req := &message.Request{}
	if ctx.Request.Body() == nil {
		return req, nil
	}

	if err := c.Unmarshal(ctx.Request.Body(), req); err != nil {
		return nil, err
	}
	req.Payload = &message.Payload{
		Header: map[string]string{"Content-Type": ct},
	}
	req.Payload.Body = ctx.Request.Body()

	return req, nil
}

func (p *Server) parseToResponse(ctx *routing.Context, msg *message.Response) error {

	ct := string(ctx.Request.Header.Peek("Content-Type"))
	c := codec.Select(ct)
	if c == nil {
		return errcode.Newf("not supported content type: %s", ct)
	}

	if msg == nil {
		msg = &message.Response{}
	}

	bs, err := c.Marshal(msg)
	if err != nil {
		return err
	}

	ctx.SetContentType(codec.ContentTypeJson)
	ctx.SetStatusCode(http.StatusOK)
	fmt.Println(bs)
	ctx.SetBody(bs)

	return nil
}
