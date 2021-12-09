package component

import (
	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis/common.v0/config"
	"trellis.tech/trellis/common.v0/logger"
)

type Component interface {
	lifecycle.LifeCycle

	Route(*message.Payload) (interface{}, error)
}

// Option 处理参数函数
type Option func(*Options)

// Options 参数对象
type Options struct {
	Logger logger.Logger
	Config config.Config
}

// Config 注入配置
func Config(c config.Config) Option {
	return func(p *Options) {
		p.Config = c
	}
}

// Logger 日志记录
func Logger(l logger.Logger) Option {
	return func(p *Options) {
		p.Logger = l
	}
}

type NewComponentFunc func(...Option) (Component, error)

func RegisterNewComponentFunc(s *service.Service, newFunc NewComponentFunc) {
	if err := m.RegisterNewComponentFunc(s, newFunc); err != nil {
		panic(err)
	}
}

func RegisterComponent(s *service.Service, comp Component) {
	if err := m.RegisterComponent(s, comp); err != nil {
		panic(err)
	}
}
