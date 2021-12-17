package registry

import (
	"flag"

	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v0/clients/etcd"
)

// NewRegistryFunc new registry function
type NewRegistryFunc func(...Option) (Registry, error)

// Registry The registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, ...}
type Registry interface {
	ID() string
	String() string

	ProcessService

	lifecycle.LifeCycle
}

// ProcessService Process register service
type ProcessService interface {
	Register(*service.ServiceNode) error
	Deregister(*service.ServiceNode) error
	Watch(*service.Service) (Watcher, error)
}

type Config struct {
	RegisterType   RegisterType `yaml:"register_type" json:"register_type"`
	RegisterPrefix string       `yaml:"register_prefix" json:"register_prefix"`
	RetryTimes     int          `yaml:"retry_times" json:"retry_times"`

	RegisterServices `yaml:",inline" json:",inline"`
	WatchServices    []*WatchService `yaml:"watch_services" json:"watch_services"`

	ETCDConfig etcd.Config `yaml:"etcd_config" json:"etcd_config"`
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *Config) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	registryType := f.Int(prefix+"registry.register_type", int(RegisterType_memory),
		"The register type of router. 1: etcd, default: memory.")
	cfg.RegisterType = RegisterType(*registryType)
	f.StringVar(&cfg.RegisterPrefix, prefix+"registry.register_prefix", "/", "The register prefix.")
	f.IntVar(&cfg.RetryTimes, prefix+"registry.retry_times", 0, "The register retry times of nodes.")
	//f.Var(&cfg.Heartbeat, prefix+"registry.heartbeat", "The register heartbeat.")
	//f.Var(&cfg.TTL, prefix+"registry.ttl", "The register ttl.")

	cfg.ETCDConfig.ParseFlagsWithPrefix(prefix, f)
}
