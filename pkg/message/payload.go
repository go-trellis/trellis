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
	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis/common.v1/errcode"
	"trellis.tech/trellis/common.v1/json"
)

func (x *Payload) ContentType() string {
	header := x.GetHeader()
	if header == nil {
		return ""
	}
	return header["Content-Type"]
}

func (x *Payload) ToObject(model interface{}) error {
	if x.GetBody() == nil {
		return nil
	}
	c, err := x.getCodec()
	if err != nil {
		return err
	}
	return c.Unmarshal(x.GetBody(), model)
}

func (x *Payload) Set(k, v string) {
	if x == nil {
		x = &Payload{}
	}
	if x.Header == nil {
		x.Header = make(map[string]string)
	}
	x.Header[k] = v
}

func (x *Payload) Get(k string) string {
	if x == nil || x.Header == nil {
		return ""
	}
	return x.Header[k]
}

func (x *Payload) SetBody(model interface{}) error {
	switch t := model.(type) {
	case []byte:
		x.Body = t
		return nil
	default:
		c, err := x.getCodec()
		if err != nil {
			return err
		}

		x.Body, err = c.Marshal(model)
		if err != nil {
			return err
		}
		return nil
	}
}

func (x *Payload) GetTraceInfo() (*mime.TraceInfo, error) {
	header := x.GetHeader()
	if header == nil {
		return nil, errcode.New("nil header")
	}

	bs, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}

	info := &mime.TraceInfo{}
	if err = json.Unmarshal(bs, info); err != nil {
		return nil, err
	}
	return info, nil
}

func (x *Payload) getCodec() (codec.Codec, error) {
	ct := x.ContentType()
	c := codec.Select(ct)
	if c == nil {
		return nil, errcode.Newf("not supported content-type: %q", ct)
	}
	return c, nil
}
