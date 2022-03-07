package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/router"

	_ "trellis.tech/trellis.v1/examples/components"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/server/http"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/trellis"

	"trellis.tech/trellis/common.v1/config"
)

func main() {

	r, err := router.NewRouter(router.Config{
		RegistryConfig: registry.Config{
			RegisterType:     registry.RegisterType_memory,
			RegisterPrefix:   "/trellis",
			RegisterServices: registry.RegisterServices{},
			WatchServices:    []*registry.WatchService{},
		},
		//ETCDConfig     etcd.Config
		//Logger: logger.Noop(),
		Components: []*component.Config{&component.Config{
			Service: service.NewService("trellis", "componentb", "v1"),
			Options: config.Options{"server": "componentb"},
		}}},
	)

	if err != nil {
		panic(err)
	}

	s, err := http.NewServer(
		http.Config(&trellis.HTTPServerConfig{Address: "0.0.0.0:8000"}),
		http.Router(r),
	)
	if err != nil {
		panic(err)
	}

	s.RegisterHandlers(&trellis.HTTPHandler{
		Method:  "POST",
		Path:    "/v1",
		Uses:    []string{"use1"},
		Handler: "",
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
