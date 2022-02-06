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
	c, err := x.getCodec()
	if err != nil {
		return err
	}

	return c.Unmarshal(x.GetBody(), model)
}

func (x *Payload) Set(k, v string) {
	if x.Header == nil {
		x.Header = make(map[string]string)
	}
	x.Header[k] = v
}

func (x *Payload) Get(k string) string {
	if x == nil {
		return ""
	}
	if x.Header == nil {
		return ""
	}
	return x.Header[k]
}

func (x *Payload) SetBody(model interface{}) (err error) {
	c, err := x.getCodec()
	if err != nil {
		return err
	}

	x.Body, err = c.Marshal(model)
	return err
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
