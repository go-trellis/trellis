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

package registry

import (
	"flag"
	"path/filepath"

	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/clients/etcd"
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
	Register(*ServiceNode) error
	Deregister(*ServiceNode) error
	Watch(*service.Service) (Watcher, error)
}

type RegisterServices struct {
	RegisterServiceNodes []*ServiceNode `yaml:"register_service_nodes" json:"register_service_nodes"`
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
	registryType := f.Int(prefix+"registry.register_type", int(RegisterType_REGISTER_TYPE_MEMORY),
		"The register type of router. 1: etcd, default: memory.")
	cfg.RegisterType = RegisterType(*registryType)
	f.StringVar(&cfg.RegisterPrefix, prefix+"registry.register_prefix", "/", "The register prefix.")
	f.IntVar(&cfg.RetryTimes, prefix+"registry.retry_times", 0, "The register retry times of nodes.")
	//f.Var(&cfg.Heartbeat, prefix+"registry.heartbeat", "The register heartbeat.")
	//f.Var(&cfg.TTL, prefix+"registry.ttl", "The register ttl.")

	cfg.ETCDConfig.ParseFlagsWithPrefix(prefix+"registry.", f)
}

func (x *ServiceNode) RegisteredServiceNode(registry string) string {
	x.Service.Domain = service.CheckDomain(x.GetService().GetDomain())
	return filepath.Join(registry, x.Service.GetDomain(),
		service.ReplaceURL(x.Service.GetName()),
		service.ReplaceURL(x.Service.GetVersion()),
		service.ReplaceURL(x.Node.Value))
}
