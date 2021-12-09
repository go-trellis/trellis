package codec

import "trellis.tech/trellis/common.v0/json"

type JSON struct{}

func NewJsonCodec() (Codec, error) {
	return &JSON{}, nil
}

func (p *JSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (*JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
