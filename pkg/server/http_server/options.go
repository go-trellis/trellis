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

package http_server

import (
	"fmt"

	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/trellis"

	routing "github.com/go-trellis/fasthttp-routing"
	"trellis.tech/trellis/common.v1/errcode"
)

type Handler struct {
	Method  string
	Path    string
	Uses    []routing.Handler
	Handler routing.Handler
	Service *service.Service
}

func (p *Handler) fullpath() string {
	return queryFullpath(p.Path, p.Method)
}

func queryFullpath(url, method string) string {
	return fmt.Sprintf("%s:%s", url, method)
}

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

func (p *Server) RegisterGroupHandlers(cfgs ...*trellis.HTTPGroup) error {
	return p.registerGroupHandlers(p.fastRouter, cfgs...)
}

func (p *Server) registerGroupHandlers(router *routing.Router, cfgs ...*trellis.HTTPGroup) error {
	for _, cfg := range cfgs {
		group := router.Group(cfg.Path)

		var groupUses []routing.Handler
		for _, use := range cfg.Uses {
			uFunc, err := server.GetUseFunc(use)
			if err != nil {
				return err
			}
			groupUses = append(groupUses, uFunc)
		}

		for _, hCfg := range cfg.Handlers {
			handler, err := p.getHTTPHandler(hCfg, groupUses)
			if err != nil {
				return err
			}
			group.To(handler.Method, handler.Path, append(handler.Uses, handler.Handler)...)

			p.services[cfg.Path+handler.fullpath()] = handler.Service
		}
	}

	return nil
}

func (p *Server) RegisterHandlers(cfgs ...*trellis.HTTPHandler) error {
	return p.registerHandlers(p.fastRouter, cfgs...)
}

func (p *Server) registerHandlers(router *routing.Router, cfgs ...*trellis.HTTPHandler) error {
	for _, cfg := range cfgs {
		handler, err := p.getHTTPHandler(cfg, nil)
		if err != nil {
			return err
		}

		router.To(handler.Method, handler.Path, append(handler.Uses, handler.Handler)...)
		p.services[handler.fullpath()] = handler.Service
	}
	return nil
}

func (p *Server) getHTTPHandler(handler *trellis.HTTPHandler, groupUses []routing.Handler) (*Handler, error) {
	if handler.Service == nil {
		return nil, errcode.Newf("not set service to handler : %s", handler.Path)
	}
	h := &Handler{
		Method:  handler.Method,
		Path:    handler.Path,
		Uses:    groupUses,
		Service: handler.Service,
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
		h.Handler = uFunc
	} else {
		h.Handler = p.HandleHTTP
	}

	return h, nil
}
