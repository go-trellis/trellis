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

package component

import (
	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/config"
	"trellis.tech/trellis/common.v1/logger"
)

type NewComponentFunc func(*Config) (Component, error)

type Component interface {
	lifecycle.LifeCycle

	Route(topic string, payload *message.Payload) (interface{}, error)
}

type Manager interface {
	RegisterNewComponentFunc(s *service.Service, newFunc NewComponentFunc) error
	RegisterComponent(s *service.Service, component Component) error
	NewComponent(c *Config) error
	GetComponent(*service.Service) Component
	StopComponents() error
}

type Configs []*Config

type Config struct {
	Service *service.Service `yaml:"service" json:"service"`
	Options config.Options   `yaml:"options" json:"options"`

	Caller clients.Caller `yaml:"-" json:"-"`
	Logger logger.Logger  `yaml:"-" json:"-"`
}
