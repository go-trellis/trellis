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

package client

import (
	"reflect"
	"sync"

	"trellis.tech/trellis.v1/pkg/clients/grpc"
	"trellis.tech/trellis.v1/pkg/clients/http"
	"trellis.tech/trellis.v1/pkg/clients/local"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/server"

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

func New(nd *node.Node) (server.Caller, []server.CallOption, error) {
	if nd == nil {
		return nil, nil, errcode.New("nil node")
	}
	switch nd.Protocol {
	case node.Protocol_PROTOCOL_LOCAL:
		return local.NewClient()
	case node.Protocol_PROTOCOL_GRPC:
		var (
			c   *grpc.Client
			err error
		)
		ci, ok := getClient(nd.String())
		if !ok {
			c, err = grpc.NewClient(nd)
			if err != nil {
				return nil, nil, err
			}
		} else {
			c, ok = ci.(*grpc.Client)
			if !ok {
				return nil, nil, errcode.Newf("client is not grpc client: %s", reflect.TypeOf(ci).String())
			}
		}
		if c.Pool != nil {
			setClient(nd.String(), c)
		}
		// todo call options
		return c, nil, nil
	case node.Protocol_PROTOCOL_HTTP:
		var (
			c   *http.Client
			err error
		)
		ci, ok := getClient(nd.String())
		if ok {
			c, ok = ci.(*http.Client)
			if !ok {
				return nil, nil, errcode.Newf("client is not http_server client: %s", reflect.TypeOf(ci).String())
			}
			return c, nil, nil
		} else {
			c, err = http.NewClient(nd)
			if err != nil {
				return nil, nil, err
			}
		}
		setClient(nd.String(), c)
		return c, nil, nil
		//case node.Protocol_QUIC:
		//	return quic.NewClient(n)
	}

	return nil, nil, errcode.Newf("not supported node protocol: %d, %s", nd.Protocol, nd.Protocol.String())
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
