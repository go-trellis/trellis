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

package components

import (
	"context"
	"fmt"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/service"
)

func init() {
	component.RegisterNewComponentFunc(
		service.NewService("trellis", "componenta", "v1"), NewComponentA)
}

type ComponentA struct {
	conf *component.Config
}

func NewComponentA(c *component.Config) (component.Component, error) {
	return &ComponentA{conf: c}, nil
}

func (p *ComponentA) Start() error {
	println("I am test component a start")
	return nil
}

func (p *ComponentA) Stop() error {
	println("I am test component a stop")
	return nil
}

type ReqComponentA struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type RespComponentA struct {
	Message string `json:"message"`
}

func (p *ComponentA) Route(topic string, msg *message.Payload) (interface{}, error) {
	fmt.Println("topic", topic, msg)
	fmt.Println(msg.GetTraceInfo())
	var req *ReqComponentA
	err := msg.ToObject(req)
	if err != nil {
		panic(err)
		return nil, err
	}

	switch topic {
	case "grpc":
		fmt.Println("I am test component route, topic: grpc", req)
		sTopic := service.NewServiceWithTopic("/trellis", "componentb", "v1", "test")

		return p.conf.Caller.Call(context.Background(), &message.Request{Service: sTopic, Payload: msg})
	default:
		fmt.Println("I am test component route, topic: default", req)

		return &message.Response{
			Code: 401,
			Msg:  "not found topic",
			Payload: &message.Payload{
				Header: map[string]string{"message": fmt.Sprintf("Hello: %s", req.Name)},
				Body:   []byte("say hello"),
			},
		}, nil
	}
}
