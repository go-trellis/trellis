package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trellis.tech/trellis/common.v1/config"

	_ "trellis.tech/trellis.v1/examples/components"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server/grpc"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/trellis"
	"trellis.tech/trellis/common.v1/clients/etcd"
)

var (
	srv = "0.0.0.0:8001"
)

// TODO example
func main() {
	flag.StringVar(&srv, "srv", "0.0.0.0:8001", "server address")
	flag.Parse()

	options := config.Options{"server": srv}
	s, err := grpc.NewServer(trellis.ServerConfig{
		Address: srv,
		RouterConfig: router.Config{
			RegistryConfig: registry.Config{
				RegisterType:   registry.RegisterType_etcd,
				RegisterPrefix: "/trellis",
				RegisterServices: registry.RegisterServices{
					RegisterServiceNodes: []*service.ServiceNode{
						&service.ServiceNode{
							Service: service.NewService("trellis", "componentb", "v1"),
							Node: &node.Node{
								BaseNode: node.BaseNode{
									Weight:    1024,
									Value:     srv,
									TTL:       uint64(time.Second * 10),
									Heartbeat: uint64(time.Second * 5),
									Protocol:  node.Protocol_GRPC,
								},
							},
						},
					},
				},
				WatchServices: []*registry.WatchService{},
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
				Service: service.NewService("trellis", "componentb", "v1"),
				Options: options.ToConfig(),
			},
		},
	})
	if err != nil {
		panic(err)
	}

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
