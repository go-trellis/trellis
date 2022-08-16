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

package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"trellis.tech/trellis.v1/pkg/trellis"
)

func InitTracer(name string, cfg *trellis.TracingConfig) (io.Closer, error) {
	if cfg == nil {
		return nil, nil
	}

	tracer, closer, err := NewJeagerTracer(name, cfg)
	if err != nil {
		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
