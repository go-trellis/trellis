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

package routes

import (
	"trellis.tech/trellis.v0/service/component"
	"trellis.tech/trellis/common.v0/logger"
)

type Option func(*Options)

type Options struct {
	manager component.Manager
	logger  logger.Logger
}

func CompManager(m component.Manager) Option {
	return func(o *Options) {
		o.manager = m
	}
}

func Logger(logger logger.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}
