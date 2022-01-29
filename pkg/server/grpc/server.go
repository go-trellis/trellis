package grpc

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc/peer"

	"github.com/google/uuid"
	"trellis.tech/trellis.v1/pkg/mime"

	"trellis.tech/trellis.v1/pkg/clients/client"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/trellis"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	ServerConfig trellis.ServerConfig

	rpcServer *grpc.Server

	routes router.Router
}

func (p *Server) Call(ctx context.Context, msg *message.Request) (*message.Response, error) {

	serviceNode, ok := p.routes.GetServiceNode(msg.GetService(), msg.String())
	if !ok {
		// TODO warn Log
	}

	c, err := client.New(serviceNode)
	if err != nil {
		return nil, err
	}
	// TODO Options

	if msg.GetPayload().GetHeader() == nil {
		msg.GetPayload().Header = map[string]string{}
	}

	ip, _ := peer.FromContext(ctx)

	if ip != nil {
		if msg.GetPayload().Header[mime.HeaderKeyClientIP] == "" {
			msg.GetPayload().Header[mime.HeaderKeyClientIP] = ip.Addr.String()
		}
		msg.GetPayload().Header[mime.HeaderKeyRequestIP] = ip.Addr.String()
	}

	msg.GetPayload().Header[mime.HeaderKeyRequestID] = uuid.New().String()

	if msg.GetPayload().GetHeader()[mime.HeaderKeyTraceID] == "" {
		msg.GetPayload().Header[mime.HeaderKeyTraceID] = uuid.New().String()
	}

	msg.GetPayload().Set(mime.HeaderKeyRequestID, uuid.New().String())
	return c.Call(ctx, msg)
}

func NewServer(conf trellis.ServerConfig) (*Server, error) {
	s := &Server{
		ServerConfig: conf,

		routes: router.NewRouter(conf.RouterConfig),
	}

	var sopts []grpc.ServerOption

	if conf.GrpcServerConfig.KeepaliveTime > 0 {
		sopts = append(sopts, grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    time.Duration(conf.GrpcServerConfig.KeepaliveTime),
			Timeout: time.Duration(conf.GrpcServerConfig.KeepaliveTimeout),
		}))
	}

	if conf.GrpcServerConfig.ConnectionTimeout > 0 {
		sopts = append(sopts, grpc.ConnectionTimeout(time.Duration(conf.GrpcServerConfig.ConnectionTimeout)))
	}

	if conf.GrpcServerConfig.NumStreamWorkers > 0 {
		sopts = append(sopts, grpc.NumStreamWorkers(conf.GrpcServerConfig.NumStreamWorkers))
	}

	if conf.EnableTLS {
		tls, err := conf.TLSConfig.GetTLSConfig()
		if err != nil {
			return nil, err
		}
		sopts = append(sopts, grpc.Creds(credentials.NewTLS(tls)))
	}

	s.rpcServer = grpc.NewServer(sopts...)

	if err := s.routes.Start(); err != nil {
		return nil, err
	}

	server.RegisterTrellisServer(s.rpcServer, s)
	return s, nil
}

func (p *Server) Start() error {

	// TODO config to new component
	for _, comp := range p.ServerConfig.Components {
		if err := router.NewComponent(comp); err != nil {
			return err
		}
	}

	listen, err := net.Listen("tcp", p.ServerConfig.Address)
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
	if err := router.StopComponents(); err != nil {
		// TODO log
	}
	//p.compManager.
	if err := p.routes.Stop(); err != nil {
		// TODO log
	}
	p.rpcServer.Stop()
	return nil
}
