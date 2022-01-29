package main

import (
	"context"
	"encoding/json"
	"fmt"

	"trellis.tech/trellis.v1/pkg/clients/http"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/service"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	c, err := http.NewClient(&node.Node{
		BaseNode: node.BaseNode{
			Value: "http://127.0.0.1:8000/v1",
		},
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

	fmt.Println(resp)

	r := &Response{}
	if err := json.Unmarshal(resp.GetPayload().GetBody(), r); err != nil {
		panic(err)
	}
	fmt.Println(r)
}
