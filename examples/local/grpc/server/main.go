package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "trellis.tech/trellis.v1/examples/components"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server/grpc"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/trellis"
)

func main() {
	s, err := grpc.NewServer(trellis.ServerConfig{
		Address: "0.0.0.0:8000",
		RouterConfig: router.Config{
			RegistryConfig: registry.Config{
				RegisterType:     registry.RegisterType_memory,
				RegisterPrefix:   "/trellis",
				RegisterServices: registry.RegisterServices{},
				WatchServices:    []*registry.WatchService{},
			},
			//ETCDConfig     etcd.Config
			//Logger: logger.Noop(),
		},
		Components: []*component.Config{
			&component.Config{
				Service: service.NewService("trellis", "componentb", "v1"),
			},
		},
	})
	if err != nil {
		panic(err)
	}

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
