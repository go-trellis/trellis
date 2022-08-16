/*
Copyright © 2022 Henry Huang <hhh@rutcode.com>
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

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
	"trellis.tech/trellis.v1/pkg/server/http_server"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/trellis"

	"trellis.tech/trellis/common.v1/config"
)

func main() {

	r, err := router.NewRouter(router.Config{
		RegistryConfig: registry.Config{
			RegisterType:     registry.RegisterType_REGISTER_TYPE_MEMORY,
			RegisterPrefix:   "/trellis",
			RegisterServices: registry.RegisterServices{},
			WatchServices:    []*registry.WatchService{},
		},
		//ETCDConfig     etcd.Config
		//Logger: logger.Noop(),
		Components: []*component.Config{
			{
				Service: service.NewService("trellis", "componentb", "v1"),
				Options: config.Options{"server": "componentb"},
			},
			{
				Service: service.NewService("trellis", "componenta", "v1"),
				Options: config.Options{"server": "componenta"},
			},
		}},
	)

	if err != nil {
		panic(err)
	}

	s, err := http_server.NewServer(
		http_server.Config(&trellis.HTTPServerConfig{Address: "0.0.0.0:8000"}),
		http_server.Router(r),
	)
	if err != nil {
		panic(err)
	}

	err = s.RegisterHandlers(&trellis.HTTPHandler{
		Method:  "POST",
		Path:    "/v1",
		Uses:    []string{"use1"},
		Handler: "",
		Service: service.NewServiceWithTopic("trellis", "componenta", "v1", "grpc"),
	})
	if err != nil {
		panic(err)
	}

	if err = s.Start(); err != nil {
		log.Fatalln(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Kill, os.Interrupt, syscall.SIGQUIT)
	<-ch

	if err = s.Stop(); err != nil {
		log.Fatalln(err)
	}
}
