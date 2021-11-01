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

package cmd

import (
	"trellis.tech/trellis.v0/routes"
	"trellis.tech/trellis.v0/sd/etcd"
	"trellis.tech/trellis.v0/sd/memory"
	"trellis.tech/trellis.v0/service"
	"trellis.tech/trellis.v0/service/component"
	"trellis.tech/trellis.v0/service/registry"
)

var (
	// DefaultNewRegistryFuncs new registry functions
	DefaultNewRegistryFuncs = map[service.RegisterType]registry.NewRegistryFunc{
		service.RegisterType_memory: memory.NewRegistry,
		service.RegisterType_etcd:   etcd.NewRegistry,
	}

	// DefaultHiddenVersions hidden versions
	DefaultHiddenVersions = []string{"0", "0.0", "0.0.0", "v0", "v0.0", "v0.0.0"}

	// DefaultCompManager default components manager
	DefaultCompManager = routes.NewCompManager()
)

// RegisterComponentFunc regist component funciton into default local route
func RegisterComponentFunc(service *service.Service, fn component.NewComponentFunc) {
	DefaultCompManager.RegisterComponentFunc(service, fn)
}
