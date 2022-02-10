package grpc

import (
	"log"
	"net"
	"net/http"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/trellis"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	name string

	conf *trellis.GrpcServerConfig

	rpcServer *grpc.Server

	router router.Router

	tracing bool
}

func (p *Server) Start() error {

	listen, err := net.Listen("tcp", p.conf.Address)
	if err != nil {
		return err
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
	if err := component.StopComponents(); err != nil {
		// TODO log
	}
	//p.compManager.
	if err := p.router.Stop(); err != nil {
		// TODO log
	}
	p.rpcServer.Stop()
	return nil
}

func NewServer(opts ...Option) (*Server, error) {
	s := &Server{}

	for _, o := range opts {
		o(s)
	}

	var sopts []grpc.ServerOption

	if s.conf.KeepaliveTime > 0 {
		sopts = append(sopts, grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    s.conf.KeepaliveTime,
			Timeout: s.conf.KeepaliveTimeout,
		}))
	}

	if s.conf.ConnectionTimeout > 0 {
		sopts = append(sopts, grpc.ConnectionTimeout(s.conf.ConnectionTimeout))
	}

	if s.conf.NumStreamWorkers > 0 {
		sopts = append(sopts, grpc.NumStreamWorkers(s.conf.NumStreamWorkers))
	}

	if s.conf.EnableTLS {
		tls, err := s.conf.TLSConfig.GetTLSConfig()
		if err != nil {
			return nil, err
		}
		sopts = append(sopts, grpc.Creds(credentials.NewTLS(tls)))
	}

	if s.tracing {
		sopts = append(sopts, grpc_middleware.WithUnaryServerChain(
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads()),
		))

		sopts = append(sopts, grpc_middleware.WithStreamServerChain(
			otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads()),
		))
	}

	s.rpcServer = grpc.NewServer(sopts...)

	if err := s.router.Start(); err != nil {
		return nil, err
	}

	server.RegisterTrellisServer(s.rpcServer, s.router)
	return s, nil
}
