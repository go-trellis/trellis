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

package message

import (
	"reflect"

	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis/common.v1/json"
)

type Option func(*Options)

type Options struct {
	Code int64
	Err  error
}

func Code(code int64) Option {
	return func(o *Options) {
		o.Code = code
	}
}

func Error(err error) Option {
	return func(o *Options) {
		o.Err = err
	}
}

func NewResponse(data interface{}, opts ...Option) *Response {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	resp := &Response{
		Code: options.Code,
		Payload: &Payload{
			Header: make(map[string]string),
		},
	}

	if options.Err != nil {
		resp.Msg = options.Err.Error()
	}

	if data == nil {
		return resp
	}

	var respData interface{}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr, reflect.Struct:
		respData = data
	default:
		respData = map[string]interface{}{"data": data}
	}
	bs, err := json.Marshal(respData)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp
	}
	resp.Payload.Header[mime.HeaderKeyContentType] = mime.ContentTypeJson
	resp.Payload.Body = bs
	return resp
}

func (m *Response) SetPayload(payload *Payload) {
	m.Payload = payload
}
