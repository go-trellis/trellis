package message

import (
	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis/common.v0/errcode"
)

func (m *Payload) ContentType() string {
	header := m.GetHeader()
	if header == nil {
		return ""
	}
	return header["Content-Type"]
}

func (m *Payload) ToObject(model interface{}) error {
	if m == nil {
		return nil
	}

	ct := m.ContentType()
	c := codec.Select(ct)
	if c == nil {
		return errcode.Newf("not supported content-type: %q", ct)
	}
	return c.Unmarshal(m.GetBody(), model)
}

func (m *Payload) Set(k, v string) {
	if m.Header == nil {
		m.Header = make(map[string]string)
	}
	m.Header[k] = v
}

func (m *Payload) Get(k string) string {
	if m == nil {
		return ""
	}
	if m.Header == nil {
		return ""
	}
	return m.Header[k]
}

func (m *Payload) SetBody(model interface{}) (err error) {
	if m == nil {
		return nil
	}

	ct := m.ContentType()
	c := codec.Select(ct)
	if c == nil {
		return errcode.Newf("not supported content-type: %q", ct)
	}
	m.Body, err = c.Marshal(model)
	return err
}
