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
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server/grpc_server"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/trellis"

	_ "trellis.tech/trellis.v1/examples/components"
	"trellis.tech/trellis/common.v1/clients/etcd"
	"trellis.tech/trellis/common.v1/config"
)

var (
	srv = "0.0.0.0:8001"
)

// TODO example
func main() {
	flag.StringVar(&srv, "srv", "0.0.0.0:8001", "server address")
	flag.Parse()

	r, err := router.NewRouter(router.Config{
		RegistryConfig: registry.Config{
			RegisterType:   registry.RegisterType_REGISTER_TYPE_ETCD,
			RegisterPrefix: "/trellis",
			RegisterServices: registry.RegisterServices{
				RegisterServiceNodes: []*registry.ServiceNode{
					&registry.ServiceNode{
						Service: service.NewService("trellis", "componentb", "v1"),
						Node: &node.Node{
							Weight:    1024,
							Value:     srv,
							TTL:       uint64(time.Second * 10),
							Heartbeat: uint64(time.Second * 5),
							Protocol:  node.Protocol_PROTOCOL_GRPC,
							Metadata:  map[string]string{"a": "b"},
						},
					},
				},
			},
			WatchServices: []*registry.WatchService{},
			ETCDConfig: etcd.Config{

				//Endpoints: []string{"127.0.0.1:2379"},
				//DialTimeout types.Duration   `yaml:"dial_timeout" json:"dial_timeout"`
				MaxRetries: 10,
				//EnableTLS   bool             `yaml:"tls_enabled" json:"enable_tls"`
				//TLS         commonTls.Config `yaml:",inline"`
				//Username    string           `yaml:"username" json:"username"`
				//Password    types.Secret     `yaml:"password" json:"password"`
			},
		},

		Components: []*component.Config{
			&component.Config{
				Service: service.NewService("trellis", "componentb", "v1"),
				Options: config.Options{"server": srv},
			},
		},
		//Logger: logger.Noop(),
	})
	if err != nil {
		panic(err)
	}

	s, err := grpc_server.NewServer(
		grpc_server.Config(&trellis.GrpcServerConfig{Address: srv}),
		grpc_server.Router(r))
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
