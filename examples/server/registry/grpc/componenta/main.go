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
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/component"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/server/http_server"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis.v1/pkg/trellis"

	routing "github.com/go-trellis/fasthttp-routing"
	_ "trellis.tech/trellis.v1/examples/components"
	"trellis.tech/trellis/common.v1/clients/etcd"
	"trellis.tech/trellis/common.v1/crypto/tls"
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

// TODO example
func main() {

	r, err := router.NewRouter(router.Config{
		RegistryConfig: registry.Config{
			RegisterType:     registry.RegisterType_REGISTER_TYPE_ETCD,
			RegisterPrefix:   "/trellis",
			RegisterServices: registry.RegisterServices{},
			WatchServices: []*registry.WatchService{
				&registry.WatchService{
					Service:  service.NewService("trellis", "componentb", "v1"),
					NodeType: node.NodeType_NODE_TYPE_CONSISTENT,
					Metadata: &registry.WatchServiceMetadata{
						ClientConfig: &clients.Config{
							GrpcPool: &clients.GrpcPoolConfig{
								Enable:      true,
								InitialCap:  10,
								MaxCap:      50,
								MaxIdle:     50,
								IdleTimeout: 10 * time.Second,
							},
							// 客户端如果没有在一定时间内使用，那么会释放链接
							GrpcKeepalive: &clients.GrpcKeepaliveConfig{
								Time:    5 * time.Second,
								Timeout: time.Second,

								PermitWithoutStream: true,
							},
							TlsEnable: false,
							TlsConfig: &tls.Config{
								CertPath:           "",
								KeyPath:            "",
								CAPath:             "",
								ServerName:         "",
								InsecureSkipVerify: true,
							},
						},
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
		},
		//Logger: logger.Noop(),
		Components: component.Configs{
			&component.Config{Service: service.NewService("trellis", "componenta", "v1")},
		},
	})
	if err != nil {
		panic(err)
	}

	s, err := http_server.NewServer(
		http_server.Config(&trellis.HTTPServerConfig{Address: "0.0.0.0:8000", IsGateway: true}),
		http_server.Router(r),
	)
	if err != nil {
		panic(err)
	}

	if err := s.RegisterHandlers(&trellis.HTTPHandler{
		Method:  "POST",
		Path:    "/v1",
		Uses:    []string{"use1"},
		Handler: "",
		Service: service.NewServiceWithTopic("trellis", "componenta", "v1", "grpc"),
	}); err != nil {
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
