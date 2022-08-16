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

package router

import (
	"flag"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/logger"
)

type Router interface {
	lifecycle.LifeCycle

	clients.Caller

	Register(s *registry.ServiceNode) error
	Deregister(s *registry.ServiceNode) error
	Watch(s *registry.WatchService) error

	GetServiceNode(s *service.Service, keys ...string) (*node.Node, bool)
}

type Config struct {
	RegistryConfig registry.Config   `yaml:"registry_config" json:"registry_config"`
	Components     component.Configs `yaml:"components" json:"components"`
}

// ParseFlagsWithPrefix adds the flags required to config this to the given FlagSet.
func (cfg *Config) ParseFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	cfg.RegistryConfig.ParseFlagsWithPrefix(prefix+"router.", f)
}

func NewRouter(cfg Config) (Router, error) {
	r := &routes{
		conf:         cfg,
		Logger:       logger.Noop(), // todo logger
		nodeManagers: make(map[string]node.Manager),
	}

	for _, compCfg := range cfg.Components {
		if compCfg.Caller == nil {
			compCfg.Caller = r
		}
		if err := component.NewComponent(compCfg); err != nil {
			return nil, err
		}
	}

	return r, nil
}
