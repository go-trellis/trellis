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
	fmt.Println("topic", *msg)
	req := ReqComponentA{}
	err := msg.ToObject(&req)
	if err != nil {
		return nil, err
	}

	switch topic {
	case "grpc":
		fmt.Println("I am test component route, topic: grpc", req)
		//c, err := grpc.NewClient(&node.Node{Value: "127.0.0.1:8001"})
		//if err != nil {
		//	return nil, err
		//}
		//return c.Call(context.Background(), &message.Request{
		//	Service: &service.Service{Domain: "trellis", Name: "componentb", Version: "v1", Topic: "test"},
		//	Payload: msg})

		return p.conf.TrellisServer.Call(context.Background(), &message.Request{
			Service: &service.Service{Domain: "trellis", Name: "componentb", Version: "v1", Topic: "test"},
			Payload: msg})
	default:
		fmt.Println("I am test component route, topic: default", req)

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
			Code:   401,
			ErrMsg: "not found topic",
			Payload: &message.Payload{
				Header: map[string]string{"message": fmt.Sprintf("Hello: %s", req.Name)},
				Body:   []byte("say hello"),
			},
		}, nil

		//return message.NewResponse(&TestResp{
		//	Message: fmt.Sprintf("Hello: %s", req.Name),
		//}, message.Code(401), message.Error(errcode.New("I am an error"))), nil

	}
}
