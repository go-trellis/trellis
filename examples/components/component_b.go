package components

import (
	"fmt"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/service"
)

func init() {
	router.RegisterNewComponentFunc(
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
}

type RespComponentB struct {
	Message string `json:"message"`
}

func (p *ComponentB) Route(topic string, msg *message.Payload) (interface{}, error) {

	srv := p.conf.Options.GetString("server")
	req := ReqComponentB{}
	err := msg.ToObject(&req)
	if err != nil {
		return nil, err
	}
	return &RespComponentB{
		Message: fmt.Sprintf("Hello: %s, I am component b: %s", req.Name, srv),
	}, nil
}
