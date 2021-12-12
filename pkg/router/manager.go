package router

import (
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/service"
)

type ComponentRouter interface {
	RegisterNewComponentFunc(s *service.Service, newFunc component.NewComponentFunc) error
	RegisterComponent(s *service.Service, component component.Component) error
	NewComponent(s *service.Service, opts ...component.Option) error
	GetComponent(*service.Service) component.Component
	StopComponents() error
}