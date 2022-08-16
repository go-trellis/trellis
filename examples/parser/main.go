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

package main

import (
	"encoding/json"
	"fmt"

	"trellis.tech/trellis.v1/examples/components"
	"trellis.tech/trellis.v1/pkg/server/http_server"
)

const data = `{"code":0,"trace_id":"","payload":"eyJtZXNzYWdlIjoiSGVsbG86IEhhaGEsIEkgYW0gY29tcG9uZW50IGI6IGhhaGEifQ=="}`

func main() {
	resp := &http_server.GatewayResponse{}

	if err := json.Unmarshal([]byte(data), resp); err != nil {
		panic(err)
	}

	payload := &components.RespComponentB{}

	if err := json.Unmarshal(resp.Payload, payload); err != nil {
		panic(err)
	}

	fmt.Println(payload)
}

// &{Hello: Haha, I am component b: 0.0.0.0:8001}
// &{Hello: Haha, I am component b: haha}
