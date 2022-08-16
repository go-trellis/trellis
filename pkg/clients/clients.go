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

package clients

import (
	"context"
	"time"

	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis/common.v1/crypto/tls"
)

type Config struct {
	GrpcPool      *GrpcPoolConfig      `yaml:"grpc_pool" json:"grpc_pool"`
	GrpcKeepalive *GrpcKeepaliveConfig `yaml:"grpc_keepalive" json:"grpc_keepalive"`

	TlsEnable bool        `yaml:"tls_enable" json:"tls_enable"`
	TlsConfig *tls.Config `yaml:"tls_config" json:"tls_config"`

	// todo http_server config
}

type GrpcPoolConfig struct {
	Enable      bool          `yaml:"enable" json:"enable"`
	InitialCap  int           `yaml:"initial_cap" json:"initial_cap"`
	MaxCap      int           `yaml:"max_cap" json:"max_cap"`
	MaxIdle     int           `yaml:"max_idle" json:"max_idle"`
	IdleTimeout time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
}

type GrpcKeepaliveConfig struct {
	Time                time.Duration `yaml:"time" json:"time"`
	Timeout             time.Duration `yaml:"timeout" json:"timeout"`
	PermitWithoutStream bool          `yaml:"permit_without_stream" json:"permit_without_stream"`
}

type Caller interface {
	Call(context.Context, *message.Request) (*message.Response, error)
}
