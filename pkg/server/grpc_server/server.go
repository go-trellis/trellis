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

package grpc_server

import (
	context "context"
	"fmt"
	"log"
	"net"
	"net/http"

	"trellis.tech/trellis.v1/pkg/router"

	"trellis.tech/trellis.v1/pkg/clients/client"
	"trellis.tech/trellis.v1/pkg/clients/local"
	"trellis.tech/trellis.v1/pkg/component"
	message "trellis.tech/trellis.v1/pkg/message"
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
		fmt.Println(err)
	}
	//p.compManager.
	if err := p.router.Stop(); err != nil {
		// TODO log
		fmt.Println(err)
	}
	p.rpcServer.Stop()
	return nil
}

func NewServer(opts ...Option) (*Server, error) {
	p := &Server{}

	for _, o := range opts {
		o(p)
	}

	var sopts []grpc.ServerOption

	if p.conf.KeepaliveTime > 0 {
		sopts = append(sopts, grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    p.conf.KeepaliveTime,
			Timeout: p.conf.KeepaliveTimeout,
		}))
	}

	if p.conf.ConnectionTimeout > 0 {
		sopts = append(sopts, grpc.ConnectionTimeout(p.conf.ConnectionTimeout))
	}

	if p.conf.NumStreamWorkers > 0 {
		sopts = append(sopts, grpc.NumStreamWorkers(p.conf.NumStreamWorkers))
	}

	if p.conf.EnableTLS {
		tls, err := p.conf.TLSConfig.GetTLSConfig()
		if err != nil {
			return nil, err
		}
		sopts = append(sopts, grpc.Creds(credentials.NewTLS(tls)))
	}

	if p.tracing {
		sopts = append(sopts, grpc_middleware.WithUnaryServerChain(
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads()),
		))

		sopts = append(sopts, grpc_middleware.WithStreamServerChain(
			otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads()),
		))
	}

	p.rpcServer = grpc.NewServer(sopts...)

	if err := p.router.Start(); err != nil {
		return nil, err
	}

	server.RegisterTrellisServer(p.rpcServer, p)
	return p, nil
}

func (p *Server) Call(ctx context.Context, msg *message.Request) (*message.Response, error) {
	var (
		c    server.Caller
		opts []server.CallOption
		err  error
	)
	serviceNode, ok := p.router.GetServiceNode(msg.GetService(), msg.String())
	if !ok {
		c, opts, err = local.NewClient()
	} else {
		c, opts, err = client.New(serviceNode)
	}
	if err != nil {
		return nil, err
	}
	return c.Call(ctx, msg, opts...)
}
