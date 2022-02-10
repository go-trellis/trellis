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

	// todo http config
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

type Client interface {
	Call(context.Context, *message.Request, ...CallOption) (*message.Response, error)
}
