package component

import (
	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v0/config"
	"trellis.tech/trellis/common.v0/logger"
)

type NewComponentFunc func(*Config) (Component, error)

type Component interface {
	lifecycle.LifeCycle

	Route(topic string, payload *message.Payload) (interface{}, error)
}

type Config struct {
	Service *service.Service `yaml:"service" json:"service"`
	Options config.Config    `yaml:"options" json:"options"`
	Logger  logger.Logger    `yaml:"-" json:"-"`
}
