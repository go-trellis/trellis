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

package grpc

import (
	"context"
	"reflect"

	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/server"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"trellis.tech/trellis/common.v1/errcode"
	"trellis.tech/trellis/common.v1/pool"
)

var _ server.Caller = (*Client)(nil)

type Client struct {
	nd   *node.Node
	Pool pool.Pool

	dialOptions []grpc.DialOption
}

func (p *Client) Call(ctx context.Context, in *message.Request, opts ...server.CallOption) (resp *message.Response, err error) {

	var (
		cc *grpc.ClientConn
	)
	if p.Pool != nil {
		c, err := p.Pool.Get()
		if err != nil {
			return nil, err
		}

		var ok bool
		cc, ok = c.(*grpc.ClientConn)
		if !ok {
			return nil, errcode.New("not found client in pool")
		}
		//nolint
		defer p.Pool.Put(cc)
	} else {
		cc, err = grpc.Dial(p.nd.Value, p.dialOptions...)
		if err != nil {
			return
		}
		defer cc.Close()
	}

	options := server.CallOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	return server.NewTrellisClient(cc).Call(ctx, in, options.GRPCCallOptions...)
}

func NewClient(nd *node.Node) (c *Client, err error) {

	if nd == nil {
		return nil, errcode.New("nil node")
	}
	client := &Client{
		nd: nd,
	}

	watchServiceConfig, ok := nd.Get("watch_service_config")

	if !ok {
		client.dialOptions = append(client.dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
		return client, nil
	}

	metadata, err := registry.ToWatchServiceMetadata(watchServiceConfig)
	if err != nil || metadata.ClientConfig == nil {
		client.dialOptions = append(client.dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
		return client, nil
	}

	if metadata.ClientConfig.GrpcKeepalive != nil {
		client.dialOptions = append(client.dialOptions,
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                metadata.ClientConfig.GrpcKeepalive.Time,
				Timeout:             metadata.ClientConfig.GrpcKeepalive.Timeout,
				PermitWithoutStream: metadata.ClientConfig.GrpcKeepalive.PermitWithoutStream,
			}),
		)
	}

	if metadata.ClientConfig.TlsEnable && metadata.ClientConfig.TlsConfig != nil {
		tlsCfg, err := metadata.ClientConfig.TlsConfig.GetTLSConfig()
		if err != nil {
			return nil, err
		}
		client.dialOptions = append(client.dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	} else {
		client.dialOptions = append(client.dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if opentracing.IsGlobalTracerRegistered() {
		client.dialOptions = append(client.dialOptions, grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads())))
		client.dialOptions = append(client.dialOptions, grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads())))
	}

	if metadata.ClientConfig.GrpcPool != nil && metadata.ClientConfig.GrpcPool.Enable {
		var opts = []pool.Option{
			pool.OptionFactory(client.poolFactory),
			pool.OptionClose(client.poolClose),
			pool.MaxIdle(metadata.ClientConfig.GrpcPool.MaxIdle),
			pool.MaxCap(metadata.ClientConfig.GrpcPool.MaxCap),
			pool.IdleTimeout(metadata.ClientConfig.GrpcPool.IdleTimeout),
			pool.InitialCap(metadata.ClientConfig.GrpcPool.InitialCap),
		}

		client.Pool, err = pool.NewPool(opts...)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (p *Client) poolFactory() (interface{}, error) {
	return grpc.Dial(p.nd.Value, p.dialOptions...)
}

func (p *Client) poolClose(x interface{}) error {
	switch t := x.(type) {
	case *grpc.ClientConn:
		return t.Close()
	case grpc.ClientConn:
		return t.Close()
	}
	return errcode.Newf("unsupported close type: %s", reflect.TypeOf(x).String())
}
