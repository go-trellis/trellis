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

package codec

import (
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis/common.v1/errcode"
)

type NewCodecFunc func() (Codec, error)

var (
	codecs        = make(map[string]NewCodecFunc)
	defaultCodecs = make(map[string]Codec)
)

func init() {
	RegisterCodec(mime.ContentTypeJson, NewJsonCodec)

	jsonCodec, _ := NewCodec(mime.ContentTypeJson)
	defaultCodecs[mime.ContentTypeJson] = jsonCodec
	defaultCodecs[mime.ContentTypeJsonBom] = jsonCodec
}

type Codec interface {
	Unmarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
}

func NewCodec(name string) (Codec, error) {
	fn, exist := codecs[name]
	if !exist {
		return nil, errcode.Newf("codec not exist %q", name)
	}
	return fn()
}

func RegisterCodec(name string, fn NewCodecFunc) {
	if len(name) == 0 {
		panic("codec name is empty")
	}

	if fn == nil {
		panic("codec fn is nil")
	}

	c, err := fn()
	if err != nil {
		panic(err)
	}

	defaultCodecs[name] = c
	codecs[name] = fn
}

func Select(contentType string) Codec {
	return defaultCodecs[contentType]
}
