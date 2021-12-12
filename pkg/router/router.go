package router

import (
	"flag"

	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis/common.v0/clients/etcd"
)

type Router interface {
	lifecycle.LifeCycle

	Register(s *service.ServiceNode) error
	Deregister(s *service.ServiceNode) error
	Watch(s *registry.WatchService) error

	GetServiceNode(s *service.Service, keys ...string) (*node.Node, bool)
}

type Config struct {
	RegistryConfig registry.Config `yaml:"registry_config" json:"registry_config"`
	ETCDConfig     etcd.Config     `yaml:"etcd_config" json:"etcd_config"`
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *Config) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	cfg.RegistryConfig.ParseFlagsWithPrefix(prefix, f)
	cfg.ETCDConfig.ParseFlagsWithPrefix(prefix, f)
}

func NewRouter(config Config) Router {
	return &routes{
		conf:         config,
		nodeManagers: make(map[string]node.Manager),
	}
}
