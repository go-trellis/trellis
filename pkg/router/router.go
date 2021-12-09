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
	Watch(s *service.Service) error

	GetServiceNode(s *service.Service, keys ...string) (*node.Node, bool)
}

type Config struct {
	RegisterType   registry.RegisterType
	NodeType       node.NodeType
	RegisterPrefix string
	ETCDConfig     etcd.Config
}

func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	cfg.ParseFlagsWithPrefix(f, "")
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *Config) ParseFlagsWithPrefix(f *flag.FlagSet, prefix string) {
	registryType := f.Int(prefix+"registry.register_type", int(registry.RegisterType_memory),
		"The register type of router. 1: etcd, default: memory.")
	nodeType := f.Uint(prefix+"registry.node_type", uint(node.NodeType_Direct), "The node type of components")
	f.StringVar(&cfg.RegisterPrefix, prefix+"registry.register_prefix", "/", "The register prefix.")

	cfg.RegisterType = registry.RegisterType(*registryType)
	cfg.NodeType = node.NodeType(*nodeType)

	cfg.ETCDConfig.ParseFlagsWithPrefix(f, prefix)
}

func NewRouter(config Config) Router {
	return &routes{conf: config}
}
