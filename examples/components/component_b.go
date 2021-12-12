package components

import (
	"fmt"

	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
)

func init() {
	router.RegisterNewComponentFunc(
		service.NewService("trellis", "componentb", "v1"), NewComponentB)
}

type ComponentB struct{}

func NewComponentB(...component.Option) (component.Component, error) {
	return &ComponentB{}, nil
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
}

type RespComponentB struct {
	Message string `json:"message"`
}

func (p *ComponentB) Route(topic string, msg *message.Payload) (interface{}, error) {
	fmt.Println("Route", topic, msg)
	req := ReqComponentB{}
	err := msg.ToObject(&req)
	if err != nil {
		return nil, err
	}
	return &RespComponentB{
		Message: fmt.Sprintf("Hello: %s, I am component b", req.Name),
	}, nil
}
