package message

import (
	"trellis.tech/trellis.v1/pkg/codec"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis/common.v1/errcode"
	"trellis.tech/trellis/common.v1/json"
)

func (m *Payload) ContentType() string {
	header := m.GetHeader()
	if header == nil {
		return ""
	}
	return header["Content-Type"]
}

func (m *Payload) ToObject(model interface{}) error {
	c, err := m.getCodec()
	if err != nil {
		return err
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

func (m *Payload) SetBody(model interface{}) error {
	switch t := model.(type) {
	case []byte:
		m.Body = t
		return nil
	default:
		c, err := m.getCodec()
		if err != nil {
			return err
		}

		m.Body, err = c.Marshal(model)
		if err != nil {
			return err
		}
		return nil
	}
}

func (m *Payload) GetTraceInfo() (*mime.TraceInfo, error) {
	header := m.GetHeader()
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

func (m *Payload) getCodec() (codec.Codec, error) {
	ct := m.ContentType()
	c := codec.Select(ct)
	if c == nil {
		return nil, errcode.Newf("not supported content-type: %q", ct)
	}
	return c, nil
}
