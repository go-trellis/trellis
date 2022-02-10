package grpc

import (
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/trellis"
)

type Option func(*Server)

func ServerName(name string) Option {
	return func(server *Server) {
		server.name = name
	}
}

func Config(c *trellis.GrpcServerConfig) Option {
	return func(server *Server) {
		server.conf = c
	}
}

func Router(r router.Router) Option {
	return func(server *Server) {
		server.router = r
	}
}

func Tracing(fs ...bool) Option {
	return func(server *Server) {
		if len(fs) == 0 {
			server.tracing = true
			return
		}
		server.tracing = fs[0]
	}
}
