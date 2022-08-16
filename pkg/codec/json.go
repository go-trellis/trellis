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
	eJson "encoding/json"
	"strconv"

	"trellis.tech/trellis/common.v1/json"
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
