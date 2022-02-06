package grpc

import (
	"context"
	"log"
	"net"
	"net/http"

	"trellis.tech/trellis.v1/pkg/component"

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
	conf *trellis.GrpcServerConfig

	rpcServer *grpc.Server

	router router.Router
}

func (p *Server) Call(ctx context.Context, msg *message.Request) (*message.Response, error) {

	serviceNode, ok := p.router.GetServiceNode(msg.GetService(), msg.String())
	if !ok {
		// TODO warn Log
	}

	c, err := client.New(serviceNode)
	if err != nil {
		return nil, err
	}

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

	msg.GetPayload().Header[mime.HeaderKeyRequestID] = uuid.NewString()

	if msg.GetPayload().GetHeader()[mime.HeaderKeyTraceID] == "" {
		msg.GetPayload().Header[mime.HeaderKeyTraceID] = uuid.NewString()
	}

	msg.GetPayload().Set(mime.HeaderKeyRequestID, uuid.NewString())
	return c.Call(ctx, msg)
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

	s.rpcServer = grpc.NewServer(sopts...)

	if err := s.router.Start(); err != nil {
		return nil, err
	}

	server.RegisterTrellisServer(s.rpcServer, s)
	return s, nil
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
