package component

import (
	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/config"
	"trellis.tech/trellis/common.v1/logger"
)

type NewComponentFunc func(*Config) (Component, error)

type Component interface {
	lifecycle.LifeCycle

	Route(topic string, payload *message.Payload) (interface{}, error)
}

type Router interface {
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

	TrellisServer server.TrellisServer `yaml:"-" json:"-"`
	Logger        logger.Logger        `yaml:"-" json:"-"`
}
