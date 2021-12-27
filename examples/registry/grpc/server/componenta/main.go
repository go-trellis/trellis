package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	routing "github.com/go-trellis/fasthttp-routing"

	_ "trellis.tech/trellis.v1/examples/components"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server/http"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/trellis"
	"trellis.tech/trellis/common.v1/clients/etcd"
)

// TODO example
func main() {
	s, err := http.NewServer(trellis.ServerConfig{
		Address: "0.0.0.0:8000",
		RouterConfig: router.Config{
			RegistryConfig: registry.Config{
				RegisterType:     registry.RegisterType_etcd,
				RegisterPrefix:   "/trellis",
				RegisterServices: registry.RegisterServices{},
				WatchServices: []*registry.WatchService{
					&registry.WatchService{
						Service:  service.NewService("trellis", "componentb", "v1"),
						NodeType: node.NodeType_Consistent,
					},
				},
			},
			ETCDConfig: etcd.Config{

				//Endpoints: []string{"127.0.0.1:2379"},
				//DialTimeout types.Duration   `yaml:"dial_timeout" json:"dial_timeout"`
				MaxRetries: 10,
				//EnableTLS   bool             `yaml:"tls_enabled" json:"enable_tls"`
				//TLS         commonTls.Config `yaml:",inline"`
				//Username    string           `yaml:"username" json:"username"`
				//Password    types.Secret     `yaml:"password" json:"password"`
			},
			//Logger: logger.Noop(),
		},
		Components: []*component.Config{
			&component.Config{
				Service: service.NewService("trellis", "componenta", "v1"),
			},
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
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Kill, os.Interrupt, syscall.SIGQUIT)
	<-ch

	if err := s.Stop(); err != nil {
		log.Fatalln(err)
	}
}
