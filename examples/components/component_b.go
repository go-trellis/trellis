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
