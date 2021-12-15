package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"trellis.tech/trellis.v1/pkg/clients/client"
	"trellis.tech/trellis.v1/pkg/clients/local"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/trellis"

	"google.golang.org/grpc"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	ServerConfig trellis.ServerConfig

	rpcServer *grpc.Server

	routes router.Router
}

func (p *Server) Call(ctx context.Context, msg *message.Request) (*message.Response, error) {

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

func NewServer(conf trellis.ServerConfig) (*Server, error) {
	s := &Server{
		ServerConfig: conf,
		rpcServer:    grpc.NewServer(),

		routes: router.NewRouter(conf.RouterConfig),
	}

	if err := s.routes.Start(); err != nil {
		return nil, err
	}

	server.RegisterTrellisServer(s.rpcServer, s)
	return s, nil
}

func (p *Server) Start() error {

	listen, err := net.Listen("tcp", p.ServerConfig.Address)
	if err != nil {
		return err
	}

	// TODO config to new component
	for _, comp := range p.ServerConfig.Components {
		if err := router.NewComponent(comp); err != nil {
			return err
		}
	}
	go func() {
		err := p.rpcServer.Serve(listen)

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
	//p.compManager.
	if err := p.routes.Stop(); err != nil {
		// TODO log
		fmt.Println(err)
	}
	p.rpcServer.Stop()
	return nil
}
