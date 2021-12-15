package trellis

import (
	"flag"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/router"

	"trellis.tech/trellis/common.v0/crypto/tls"
)

type ServerConfig struct {
	Address      string        `yaml:"address" json:"address"`
	EnableTLS    bool          `yaml:"enable_tls" json:"enable_tls"`
	TLSConfig    tls.Config    `yaml:",inline" json:",inline"`
	RouterConfig router.Config `yaml:"router_config" json:"router_config"`

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
}
