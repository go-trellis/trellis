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

package trellis

import (
	"flag"
	"time"

	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/crypto/tls"
)

type ServerConfig struct {
	ServerName string `yaml:"server_name" json:"server_name"`
	ServerType int    `yaml:"server_type" json:"server_type"`

	HTTPServerConfig HTTPServerConfig `yaml:"http_server_config" json:"http_server_config"`
	GrpcServerConfig GrpcServerConfig `yaml:"grpc_server_config" json:"grpc_server_config"`

	TracingConfig TracingConfig `yaml:"tracing_config" json:"tracing_config"`

	RouterConfig router.Config `yaml:"router_config" json:"router_config"`
}

func (cfg *ServerConfig) ParseFlags(f *flag.FlagSet) {
	cfg.ParseFlagsWithPrefix("trellis.", f)
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *ServerConfig) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.IntVar(&cfg.ServerType, prefix+"server.type", 0, "")

	//f.StringVar(&cfg.Address, prefix+"server.address", ":8000", "")
	//f.BoolVar(&cfg.EnableTLS, prefix+"server.enable_tls", false, "")
	//cfg.TLSConfig.ParseFlagsWithPrefix(prefix, f)
	cfg.RouterConfig.ParseFlagsWithPrefix(prefix, f)
	cfg.GrpcServerConfig.ParseFlagsWithPrefix(prefix, f)
	cfg.HTTPServerConfig.ParseFlagsWithPrefix(prefix, f)
}

type GrpcServerConfig struct {
	KeepaliveTime     time.Duration `yaml:"keepalive_time" json:"keepalive_time"`
	KeepaliveTimeout  time.Duration `yaml:"keepalive_timeout" json:"keepalive_timeout"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" json:"connection_timeout"`
	NumStreamWorkers  uint32        `yaml:"num_stream_workers" json:"num_stream_workers"`

	Tracing bool `yaml:"tracing" json:"tracing"`

	Address   string     `yaml:"address" json:"address"`
	EnableTLS bool       `yaml:"enable_tls" json:"enable_tls"`
	TLSConfig tls.Config `yaml:",inline" json:",inline"`
}

func (cfg *GrpcServerConfig) ParseFlags(f *flag.FlagSet) {
	cfg.ParseFlagsWithPrefix("trellis.", f)
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *GrpcServerConfig) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.DurationVar(&cfg.KeepaliveTime, prefix+"grpc.keepalive_time", 0, "")
	f.DurationVar(&cfg.KeepaliveTimeout, prefix+"grpc.keepalive_timeout", time.Second, "")
	f.DurationVar(&cfg.ConnectionTimeout, prefix+"grpc.connection_timeout", 0, "")
	f.BoolVar(&cfg.Tracing, prefix+"grpc.tracing", false, "")

	var streamWorkers uint64
	f.Uint64Var(&streamWorkers, prefix+"grpc.num_stream_workers", 0, "")
	cfg.NumStreamWorkers = uint32(streamWorkers)

	f.StringVar(&cfg.Address, prefix+"grpc.address", ":7000", "")
	f.BoolVar(&cfg.EnableTLS, prefix+"grpc.enable_tls", false, "")
	cfg.TLSConfig.ParseFlagsWithPrefix(prefix+"grpc.", f)
}

type HTTPServerConfig struct {
	Protocol string `yaml:"protocol" json:"protocol"`

	HTTP2Config HTTP2Config `yaml:",inline" json:",inline"`

	IsGateway bool `yaml:"is_gateway" json:"is_gateway"`

	Handlers []*HTTPHandler `yaml:"handlers" json:"handlers"`
	Groups   []*HTTPGroup   `yaml:"groups" json:"groups"`

	RecoverTrace     bool          `yaml:"recover_trace" json:"recover_trace"`
	Address          string        `yaml:"address" json:"address"`
	DisableKeepAlive bool          `yaml:"disable_keep_alive" json:"disable_keep_alive"`
	IdleTimeout      time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
	EnableTLS        bool          `yaml:"enable_tls" json:"enable_tls"`
	TLSConfig        tls.Config    `yaml:",inline" json:",inline"`
}

type HTTP2Config struct {
	// PingInterval is the interval at which the server will send a
	// ping message to a client.
	//
	// To disable pings set the PingInterval to a negative value.
	PingInterval time.Duration `yaml:"ping_interval" json:"ping_interval"`

	// ...
	MaxConcurrentStreams int `yaml:"max_concurrent_streams" json:"max_concurrent_streams"`

	// Debug is a flag that will allow the library to print debugging information.
	Debug bool `yaml:"debug" json:"debug"`
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *HTTP2Config) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.DurationVar(&cfg.PingInterval, prefix+"ping_interval", 0, "http2 ping interval")
	f.IntVar(&cfg.MaxConcurrentStreams, prefix+"max_concurrent_streams", 0, "http2 max_concurrent_streams")
	f.BoolVar(&cfg.Debug, prefix+"debug", false, "")
}

type HTTPHandler struct {
	Method  string           `yaml:"method" json:"method"`
	Path    string           `yaml:"path" json:"path"`
	Uses    []string         `yaml:"uses" json:"uses"`
	Handler string           `yaml:"handler" json:"handler"`
	Service *service.Service `yaml:"service" json:"service"`
}

type HTTPGroup struct {
	Path     string         `yaml:"path" json:"path"`
	Uses     []string       `yaml:"uses" json:"uses"`
	Handlers []*HTTPHandler `yaml:"handlers" json:"handlers"`
}

func (cfg *HTTPServerConfig) ParseFlags(f *flag.FlagSet) {
	cfg.ParseFlagsWithPrefix("trellis.", f)
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *HTTPServerConfig) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.StringVar(&cfg.Protocol, prefix+"http_server.protocol", "http1", "http2")
	f.StringVar(&cfg.Address, prefix+"http_server.address", ":8000", "")
	f.DurationVar(&cfg.IdleTimeout, prefix+"http_server.idle_timeout", time.Minute, "")
	f.BoolVar(&cfg.DisableKeepAlive, prefix+"http_server.disable_keep_alive", true, "")
	f.BoolVar(&cfg.EnableTLS, prefix+"http_server.enable_tls", false, "")
	f.BoolVar(&cfg.RecoverTrace, prefix+"http_server.recover_trace", false, "")
	cfg.TLSConfig.ParseFlagsWithPrefix(prefix+"http_server.", f)
	cfg.HTTP2Config.ParseFlagsWithPrefix(prefix+"http2.", f)
}

type TracingConfig struct {
	Enable bool `yaml:"enable" json:"enable"`
}
