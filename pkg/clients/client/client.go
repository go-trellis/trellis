package client

import (
	"reflect"
	"sync"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/clients/grpc"
	"trellis.tech/trellis.v1/pkg/clients/http"
	"trellis.tech/trellis.v1/pkg/clients/local"
	"trellis.tech/trellis.v1/pkg/node"

	"trellis.tech/trellis/common.v1/cache"
	"trellis.tech/trellis/common.v1/errcode"
)

var (
	clientTab   cache.TableCache
	clientMutex sync.Mutex

	TableKeySize = 10000
)

// todo too many clients, should be increase memory.
func init() {
	var err error
	clientTab, err = cache.NewTableCache("client_tab",
		//cache.OptionKeySize(TableKeySize),
		cache.OptionEvict(evictCallback()),
	)
	if err != nil {
		panic(err)
	}
}

func evictCallback() cache.EvictCallback {
	return func(key interface{}, value interface{}) {
		switch t := value.(type) {
		case *grpc.Client:
			if t.Pool != nil {
				t.Pool.Release()
			}
		case grpc.Client:
			if t.Pool != nil {
				t.Pool.Release()
			}
		}
	}
}

func New(nd *node.Node) (clients.Client, error) {

	if nd == nil {
		return local.NewClient()
	}
	switch nd.GetProtocol() {
	case node.Protocol_LOCAL:
		return local.NewClient()
	case node.Protocol_GRPC:
		var (
			c   *grpc.Client
			err error
		)
		ci, ok := getClient(nd.BaseNode.String())
		if !ok {
			c, err = grpc.NewClient(nd)
			if err != nil {
				return nil, err
			}
		} else {
			c, ok = ci.(*grpc.Client)
			if !ok {
				return nil, errcode.Newf("client is not grpc client: %s", reflect.TypeOf(ci).String())
			}
		}
		if c.Pool != nil {
			setClient(nd.BaseNode.String(), c)
		}
		return c, nil
	case node.Protocol_HTTP:
		var (
			c   *http.Client
			err error
		)
		ci, ok := getClient(nd.BaseNode.String())
		if ok {
			c, ok = ci.(*http.Client)
			if !ok {
				return nil, errcode.Newf("client is not http client: %s", reflect.TypeOf(ci).String())
			}
			return c, nil
		} else {
			c, err = http.NewClient(nd)
			if err != nil {
				return nil, err
			}
		}
		setClient(nd.BaseNode.String(), c)
		return c, nil
		//case node.Protocol_QUIC:
		//	return quic.NewClient(n)
	}

	return nil, errcode.Newf("not supported node protocol: %d, %s",
		nd.GetProtocol(), nd.GetProtocol().String())
}

func getClient(srv string) (interface{}, bool) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	vs, ok := clientTab.Lookup(srv)
	if !ok || len(vs) == 0 {
		return nil, false
	}
	return vs[0], true
}

func setClient(srv string, client interface{}) bool {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	return clientTab.Insert(srv, client)
}
