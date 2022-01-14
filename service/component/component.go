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

package component

import (
	"trellis.tech/trellis.v0/service"
	"trellis.tech/trellis.v0/service/message"
	"trellis.tech/trellis/common.v1/config"
	"trellis.tech/trellis/common.v1/logger"
)

// NewComponentFunc 服务对象生成函数申明
type NewComponentFunc func(opts ...Option) (Component, error)

// Handler handles the message function
type Handler func(message.Message) (interface{}, error)

// Middleware middlewares for next handler
type Middleware func(Handler) Handler

// Component Component
type Component interface {
	service.LifeCycle

	Route(msg message.Message) (interface{}, error)
}

// Describe description of component
type Describe struct {
	Name         string
	RegisterFunc string
	Component    Component
	Started      bool
}

// Option 处理参数函数
type Option func(*Options)

// Options 参数对象
type Options struct {
	Logger logger.Logger
	Config config.Config
	Caller message.Caller
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

// Caller remote service
func Caller(c message.Caller) Option {
	return func(p *Options) {
		p.Caller = c
	}
}
