package grpc

import (
	"context"
	"reflect"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"trellis.tech/trellis/common.v1/errcode"
	"trellis.tech/trellis/common.v1/pool"
)

var _ clients.Client = (*Client)(nil)

type Client struct {
	nd   *node.Node
	pool pool.Pool

	dialOptions []grpc.DialOption
}

func (p *Client) Call(ctx context.Context, in *message.Request) (resp *message.Response, err error) {

	var (
		cc *grpc.ClientConn
	)
	if p.pool != nil {
		c, err := p.pool.Get()
		if err != nil {
			return nil, err
		}

		var ok bool
		cc, ok = c.(*grpc.ClientConn)
		if !ok {
			return nil, errcode.New("not found client in pool")
		}
		defer p.pool.Put(cc)
	} else {

		cc, err = grpc.Dial(p.nd.Value, p.dialOptions...)
		if err != nil {
			return
		}
		defer cc.Close()
	}

	//ctx=> calloptions
	//ctx = new context

	return server.NewTrellisClient(cc).Call(ctx, in)
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

	metadata, ok := watchServiceConfig.(*registry.WatchServiceMetadata)
	if !ok || metadata.ClientConfig == nil {
		client.dialOptions = append(client.dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
		return client, nil
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

		client.pool, err = pool.NewPool(opts...)
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
