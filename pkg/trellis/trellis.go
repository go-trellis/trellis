package trellis

import (
	"flag"
	"time"

	"trellis.tech/trellis.v1/pkg/router"

	"trellis.tech/trellis/common.v1/crypto/tls"
)

type ServerConfig struct {
	Type int32 `yaml:"type" json:"type"`

	HTTPServerConfig HTTPServerConfig `yaml:"http_server_config" json:"http_server_config"`
	GrpcServerConfig GrpcServerConfig `yaml:"grpc_server_config" json:"grpc_server_config"`

	RouterConfig router.Config `yaml:"router_config" json:"router_config"`
}

func (cfg *ServerConfig) ParseFlags(f *flag.FlagSet) {
	cfg.ParseFlagsWithPrefix("trellis.", f)
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *ServerConfig) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	var serverType int
	f.IntVar(&serverType, prefix+"server.type", 0, "")
	cfg.Type = int32(serverType)
	//
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

	var streamWorkers uint64
	f.Uint64Var(&streamWorkers, prefix+"grpc.num_stream_workers", 0, "")
	cfg.NumStreamWorkers = uint32(streamWorkers)

	f.StringVar(&cfg.Address, prefix+"grpc.address", ":7000", "")
	f.BoolVar(&cfg.EnableTLS, prefix+"grpc.enable_tls", false, "")
	cfg.TLSConfig.ParseFlagsWithPrefix(prefix+"grpc.", f)
}

type HTTPServerConfig struct {
	Handlers []*HTTPHandler `yaml:"handlers" json:"handlers"`
	Groups   []*HTTPGroup   `yaml:"groups" json:"groups"`

	Address   string     `yaml:"address" json:"address"`
	EnableTLS bool       `yaml:"enable_tls" json:"enable_tls"`
	TLSConfig tls.Config `yaml:",inline" json:",inline"`
}

type HTTPHandler struct {
	Method  string   `yaml:"method" json:"method"`
	Path    string   `yaml:"path" json:"path"`
	Uses    []string `yaml:"uses" json:"uses"`
	Handler string   `yaml:"handler" json:"handler"`
}

type HTTPGroup struct {
	Path     string         `yaml:"path" json:"path"`
	Handlers []*HTTPHandler `yaml:"handlers" json:"handlers"`
}

func (cfg *HTTPServerConfig) ParseFlags(f *flag.FlagSet) {
	cfg.ParseFlagsWithPrefix("trellis.", f)
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *HTTPServerConfig) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.StringVar(&cfg.Address, prefix+"http.address", ":8000", "")
	f.BoolVar(&cfg.EnableTLS, prefix+"http.enable_tls", false, "")
	cfg.TLSConfig.ParseFlagsWithPrefix(prefix+"http.", f)
}
