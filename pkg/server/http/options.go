package http

import (
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/trellis"
)

type Option func(*Server)

func Config(c *trellis.HTTPServerConfig) Option {
	return func(server *Server) {
		server.conf = c
	}
}

func Router(r router.Router) Option {
	return func(server *Server) {
		server.router = r
	}
}
