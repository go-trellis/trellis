package http

import (
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/trellis"
)

type Option func(*Server)

func ServerName(name string) Option {
	return func(server *Server) {
		server.name = name
	}
}

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

func Tracing(fs ...bool) Option {
	return func(server *Server) {
		if len(fs) == 0 {
			server.tracing = true
			return
		}
		server.tracing = fs[0]
	}
}

func getHTTPGroupHandlers(group *trellis.HTTPGroup) ([]*Handler, error) {

	var hs []*Handler
	for _, hCfg := range group.Handlers {
		h, err := getHTTPHandler(hCfg)
		if err != nil {
			return nil, err
		}
		hs = append(hs, h)
	}

	return hs, nil
}

func getHTTPHandler(handler *trellis.HTTPHandler) (*Handler, error) {
	h := &Handler{
		Method: handler.Method,
		Path:   handler.Path,
	}

	for _, use := range handler.Uses {
		uFunc, err := server.GetUseFunc(use)
		if err != nil {
			return nil, err
		}
		h.Uses = append(h.Uses, uFunc)
	}

	if handler.Handler != "" {
		uFunc, err := server.GetUseFunc(handler.Handler)
		if err != nil {
			return nil, err
		}
		h.Uses = append(h.Uses, uFunc)
	}
	return h, nil
}
