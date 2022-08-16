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
	"context"
	"fmt"

	"trellis.tech/trellis.v1/pkg/clients/grpc"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/json"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	c, err := grpc.NewClient(&node.Node{
		Value: "127.0.0.1:8001",
	})
	if err != nil {
		panic(err)
	}
	resp, err := c.Call(context.Background(),

		&message.Request{
			Service: service.NewService("trellis", "componentb", "v1"),
			Payload: &message.Payload{
				Header: map[string]string{"Content-Type": "application/json"},
				Body:   []byte(`{"name":"haha", "age": 10}`),
			},
		})
	if err != nil {
		panic(err)
	}

	r := &Response{}
	if err := json.Unmarshal(resp.GetPayload().GetBody(), r); err != nil {
		panic(err)
	}
	fmt.Println(r)

}
