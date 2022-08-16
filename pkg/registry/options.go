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

package registry

import (
	"trellis.tech/trellis/common.v1/clients/etcd"
)

// Option initial options' functions
type Option func(*Options)

// Options new registry Options
type Options struct {
	Prefix     string
	ETCDConfig *etcd.Config
	RetryTimes int
}

func Prefix(pre string) Option {
	return func(o *Options) {
		o.Prefix = pre
	}
}

func EtcdConfig(c *etcd.Config) Option {
	return func(o *Options) {
		o.ETCDConfig = c
	}
}

func RetryTimes(times int) Option {
	return func(o *Options) {
		o.RetryTimes = times
	}
}
