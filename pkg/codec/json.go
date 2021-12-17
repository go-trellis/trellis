package codec

import (
	eJson "encoding/json"
	"strconv"

	"trellis.tech/trellis/common.v0/json"
)

type Number eJson.Number

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

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
