package main

import (
	"fmt"

	"trellis.tech/trellis.v1/pkg/server"

	routing "github.com/go-trellis/fasthttp-routing"
	"trellis.tech/trellis.v1/cmd"
	_ "trellis.tech/trellis.v1/examples/components"
)

var (
	use1 routing.Handler = func(*routing.Context) error {
		fmt.Println("I am an use handler")

		return nil
	}
)

func init() {
	server.RegisterUseFunc("use1", use1)
}

func main() {
	cmd.Run()
}
