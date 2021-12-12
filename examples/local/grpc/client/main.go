package main

import (
	"context"
	"fmt"

	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis.v1/pkg/message"

	"trellis.tech/trellis.v1/pkg/clients/grpc"
	"trellis.tech/trellis.v1/pkg/node"
)

func main() {
	c, err := grpc.NewClient(&node.Node{
		Value: "127.0.0.1:8000",
	})
	if err != nil {
		panic(err)
	}
	resp, err := c.Call(context.Background(), &message.Request{
		Service: service.NewService("trellis", "componentb", "v1"),
		Payload: &message.Payload{
			Header: map[string]string{"Content-Type": "application/json"}, Body: []byte(`{"name":"haha", "age": 10}`),
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
