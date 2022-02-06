package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"sync"
	"time"

	"trellis.tech/trellis.v1/pkg/mime"

	"trellis.tech/trellis.v1/pkg/clients/http"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/service"
)

var (
	mu  = sync.Mutex{}
	num int
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	flag.IntVar(&num, "num", 1, "request count")
	flag.Parse()

	cc := map[string]int{}
	ch := make(chan int, 100)
	for i := 0; i < num; i++ {
		ch <- i
		go func(i int) {
			r := call()
			mu.Lock()
			cc[r]++
			mu.Unlock()
			<-ch
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println(cc)
}

func call() string {
	c, err := http.NewClient(&node.Node{
		BaseNode: node.BaseNode{
			Value: "http://127.0.0.1:8000/v1",
		},
	})
	if err != nil {
		return err.Error()
	}
	s := service.NewService("trellis", "componentb", "v1")
	s.Topic = "test"
	response, err := c.Call(context.Background(),
		&message.Request{
			Service: s,
			Payload: &message.Payload{
				Header: map[string]string{"Content-Type": "application/json"},
				Body:   []byte(`{"name":"haha", "age": 10}`),
			},
		})

	if err != nil {
		return err.Error()
	}

	ct := response.GetPayload().Get(mime.HeaderKeyContentType)
	if ct == mime.ContentTypeJson {
		r := map[string]interface{}{}
		_ = json.Unmarshal(response.GetPayload().GetBody(), &r)
		return r["message"].(string)
	}
	return fmt.Sprintf("content-type err: %s", ct)
}
