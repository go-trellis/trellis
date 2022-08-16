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
	"fmt"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/service"
)

func init() {
	component.RegisterNewComponentFunc(
		service.NewService("trellis", "componentb", "v1"), NewComponentB)
}

type ComponentB struct {
	conf *component.Config
}

func NewComponentB(config *component.Config) (component.Component, error) {
	return &ComponentB{conf: config}, nil
}

func (p *ComponentB) Start() error {
	println("I am test component b start")
	return nil
}

func (p *ComponentB) Stop() error {
	println("I am test component b stop")
	return nil
}

type ReqComponentB struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type RespComponentB struct {
	Message string `json:"message"`
}

func (p *ComponentB) Route(topic string, msg *message.Payload) (interface{}, error) {
	fmt.Println("topic", topic, msg)
	fmt.Println(msg.GetTraceInfo())
	srv := p.conf.Options["server"]
	req := ReqComponentB{}
	err := msg.ToObject(&req)
	if err != nil {
		return nil, err
	}
	fmt.Println(req)
	return &RespComponentB{
		Message: fmt.Sprintf("Hello: %s, I am component b: %s", req.Name, srv),
	}, nil
}
