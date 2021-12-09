package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	routing "github.com/go-trellis/fasthttp-routing"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server/http"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis/common.v0/logger"
)

func main() {
	s, err := http.NewServer(http.Config{
		Address: "0.0.0.0:8000",
		RouterConfig: router.Config{
			RegisterType:   registry.RegisterType_memory,
			NodeType:       node.NodeType_Direct,
			RegisterPrefix: "trellis",
			//ETCDConfig     etcd.Config
			Logger: logger.Noop(),
		},
	})
	if err != nil {
		panic(err)
	}

	s.RegisterHandler(http.Handler{
		Method: "POST",
		Path:   "/v1",
		Uses: []routing.Handler{
			func(*routing.Context) error {
				fmt.Println("I am an use handler")
				return nil
			},
			//func(ctx *routing.Context) error {
			//
			//	fmt.Println("I am an error use handler")
			//	return routing.NewHTTPError(404, `{"code": 404}`)
			//},
		},
		Handler: s.HandleHTTP,
	})

	if err := s.Start(); err != nil {
		log.Fatalln(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Kill, os.Interrupt, syscall.SIGQUIT)
	<-ch

	if err := s.Stop(); err != nil {
		log.Fatalln(err)
	}
}

var _ component.Component = (*TestComponent)(nil)

func init() {
	component.RegisterComponent(&service.Service{
		Domain:  "trellis",
		Name:    "test",
		Version: "v1",
	}, &TestComponent{})
}

type TestComponent struct {
}

func (p *TestComponent) Start() error {
	println("I am test component start")
	return nil
}

func (p *TestComponent) Stop() error {
	println("I am test component stop")
	return nil
}

type TestReq struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type TestResp struct {
	Message string `json:"message"`
}

func (p *TestComponent) Route(msg *message.Payload) (interface{}, error) {
	fmt.Println("Route", msg)
	req := TestReq{}
	err := msg.ToObject(&req)
	if err != nil {
		return nil, err
	}
	fmt.Println("I am test component route", req)
	//return nil, nil
	//return nil, errcode.New("I am response an error")
	//return &TestResp{
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
