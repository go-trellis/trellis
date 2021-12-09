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

package registry

import (
	"trellis.tech/trellis/common.v0/logger"
)

// Option initial options' functions
type Option func(*Options)

// Options new registry Options
type Options struct {
	Prefix string

	Logger logger.Logger
}

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func Prefix(pre string) Option {
	return func(o *Options) {
		o.Prefix = pre
	}
}

//// DeregisterOption options' of deregistering service functions
//type DeregisterOption func(*DeregisterOptions)
//
//// DeregisterOptions deregister service Options
//type DeregisterOptions struct {
//	TTL time.Duration
//	// Other options for implementations of the interface
//	// can be stored in a context
//	Context context.Context
//}
//
//// WatchOption options' of watching service functions
//type WatchOption func(*WatchOptions)
//
//// WatchOptions watch service Options
//type WatchOptions struct {
//	Logger logger.Logger
//}
//
//func WatchLogger(l logger.Logger) WatchOption {
//	return func(w *WatchOptions) {
//		w.Logger = l
//	}
//}
