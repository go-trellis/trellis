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

package etcd

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/service"

	clientv3 "go.etcd.io/etcd/client/v3"
	"trellis.tech/trellis/common.v1/clients/etcd"
)

type etcdWatcher struct {
	registryID string

	w       clientv3.WatchChan
	client  etcd.Clientv3Facade
	timeout time.Duration

	sync.Mutex

	cancel func()
	stop   chan bool
}

func newEtcdWatcher(c etcd.Clientv3Facade, e *etcdRegistry, s *service.Service) (
	registry.Watcher, error) {

	if s.GetName() == "" {
		return nil, errors.New("service name not found")
	}

	ctx, cancel := context.WithCancel(context.Background())

	w := c.Watch(ctx, s.GetPath(e.options.Prefix), clientv3.WithPrefix(), clientv3.WithPrevKV())
	stop := make(chan bool, 1)

	return &etcdWatcher{
		registryID: e.id,
		cancel:     cancel,
		stop:       stop,
		w:          w,
		client:     c,
		timeout:    time.Second * 5, // todo config
	}, nil
}

func (p *etcdWatcher) Next() (*registry.Result, error) {
	for resp := range p.w {
		if resp.Err() != nil {
			return nil, resp.Err()
		}

		if resp.Canceled {
			return nil, errors.New("could not get next")
		}

		for _, ev := range resp.Events {
			s := decode(ev.Kv.Value)
			var typ service.EventType

			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					typ = service.EventType_EVENT_TYPE_CREATE
				} else if ev.IsModify() {
					typ = service.EventType_EVENT_TYPE_UPDATE
				}
			case clientv3.EventTypeDelete:
				typ = service.EventType_EVENT_TYPE_DELETE

				// get service from prevKv
				s = decode(ev.PrevKv.Value)
			}

			fmt.Println(s)

			if s == nil {
				continue
			}

			return &registry.Result{
				Id:          p.registryID,
				EventType:   typ,
				Timestamp:   time.Now().UnixNano(),
				ServiceNode: s,
			}, nil
		}
	}
	return nil, errors.New("could not get next")
}

func (p *etcdWatcher) Stop() {
	p.Lock()
	defer p.Unlock()

	select {
	case <-p.stop:
		return
	default:
		close(p.stop)
		p.cancel()
		p.client.Close()
	}
}
