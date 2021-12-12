/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

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

	RegisterServices `yaml:",inline" json:",inline"`
	WatchServices    []*WatchService `yaml:"watch_servics" json:"watch_servics"`

	ETCDConfig etcd.Config `yaml:"etcd_config" json:"etcd_config"`
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *Config) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	registryType := f.Int(prefix+"registry.register_type", int(RegisterType_memory),
		"The register type of router. 1: etcd, default: memory.")
	cfg.RegisterType = RegisterType(*registryType)
	f.StringVar(&cfg.RegisterPrefix, prefix+"registry.register_prefix", "/", "The register prefix.")

	cfg.ETCDConfig.ParseFlagsWithPrefix(prefix, f)
}
