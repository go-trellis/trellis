package main

import (
	"trellis.tech/trellis.v1/pkg/clients/http"
	"trellis.tech/trellis.v1/pkg/node"
)

func main() {
	http.NewClient(&node.Node{
		Value: "http://127.0.0.1:8000/",
	})
}
