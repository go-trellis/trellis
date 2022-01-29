package trellis

import (
	"flag"
	"time"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/router"

	"trellis.tech/trellis/common.v1/crypto/tls"
)

type ServerConfig struct {
	Address      string        `yaml:"address" json:"address"`
	EnableTLS    bool          `yaml:"enable_tls" json:"enable_tls"`
	TLSConfig    tls.Config    `yaml:",inline" json:",inline"`
	RouterConfig router.Config `yaml:"router_config" json:"router_config"`

	GrpcServerConfig GrpcServerConfig `yaml:"grpc_server_config" json:"grpc_server_config"`

	Components []*component.Config `yaml:"components" json:"components"`
}

func (cfg *ServerConfig) ParseFlags(f *flag.FlagSet) {
	cfg.ParseFlagsWithPrefix("trellis.", f)
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *ServerConfig) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.StringVar(&cfg.Address, prefix+"server.address", ":8000", "")
	f.BoolVar(&cfg.EnableTLS, prefix+"server.enable_tls", false, "")
	cfg.TLSConfig.ParseFlagsWithPrefix(prefix, f)
	cfg.RouterConfig.ParseFlagsWithPrefix(prefix, f)
	cfg.GrpcServerConfig.ParseFlagsWithPrefix(prefix, f)
}

type GrpcServerConfig struct {
	KeepaliveTime     time.Duration `yaml:"keepalive_time" json:"keepalive_time"`
	KeepaliveTimeout  time.Duration `yaml:"keepalive_timeout" json:"keepalive_timeout"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" json:"connection_timeout"`
	NumStreamWorkers  uint32        `yaml:"num_stream_workers" json:"num_stream_workers"`
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
}
