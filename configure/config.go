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

package configure

import "trellis.tech/trellis/common.v1/logger"

type Configure struct {
	Project Project `json:"project" yaml:"project"`
}

type Project struct {
	Logger     logger.LogConfig     `json:"logger" yaml:"logger"`
	Registries map[string]*Registry `json:"registries" yaml:"registries"`
	Services   map[string]*Service  `json:"services" yaml:"services"`
}
