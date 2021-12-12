package components

import (
	"fmt"

	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
)

func init() {
	router.RegisterComponent(
		service.NewService("trellis", "componenta", "v1"), &ComponentA{})
}

type ComponentA struct{}

func NewComponentA(...component.Option) (component.Component, error) {
	return &ComponentA{}, nil
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
	fmt.Println("Route", topic, msg)
	req := ReqComponentA{}
	err := msg.ToObject(&req)
	if err != nil {
		return nil, err
	}
	fmt.Println("I am test component route", req)
	//return nil, nil
	//return nil, errcode.New("I am response an error")
	//return &RespComponentA{
	//	Message: fmt.Sprintf("Hello: %s", req.Name),
	//}, nil
	//return &message.Payload{
	//	Header: map[string]string{"message": fmt.Sprintf("Hello: %s", req.Name)},
	//	Body:   []byte("say hello"),
	//}, nil

	//return message.NewResponse(&TestResp{
	//	Message: fmt.Sprintf("Hello: %s", req.Name),
	//}, message.Code(401)), nil

	return &message.Response{
		Code: 401,
		Payload: &message.Payload{
			Header: map[string]string{"message": fmt.Sprintf("Hello: %s", req.Name)},
			Body:   []byte("say hello"),
		},
	}, nil

	//return message.NewResponse(&TestResp{
	//	Message: fmt.Sprintf("Hello: %s", req.Name),
	//}, message.Code(401), message.Error(errcode.New("I am an error"))), nil
}
