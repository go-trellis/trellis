package router

import (
	"flag"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/logger"
)

type Router interface {
	lifecycle.LifeCycle

	Register(s *service.Node) error
	Deregister(s *service.Node) error
	Watch(s *registry.WatchService) error

	GetServiceNode(s *service.Service, keys ...string) (*node.Node, bool)

	server.Caller
}

type Config struct {
	RegistryConfig registry.Config   `yaml:"registry_config" json:"registry_config"`
	Components     component.Configs `yaml:"components" json:"components"`
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *Config) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	cfg.RegistryConfig.ParseFlagsWithPrefix(prefix+"router.", f)
}

func NewRouter(c Config) (Router, error) {
	r := &routes{
		conf:         c,
		Logger:       logger.Noop(), // todo logger
		nodeManagers: make(map[string]node.Manager),
	}

	for _, compCfg := range c.Components {
		if compCfg.Caller == nil {
			compCfg.Caller = r
		}

		if err := component.NewComponent(compCfg); err != nil {
			return nil, err
		}
	}

	return r, nil
}
