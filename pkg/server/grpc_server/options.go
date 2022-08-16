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
