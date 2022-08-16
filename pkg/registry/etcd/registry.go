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
	"fmt"
	"sync"
	"time"

	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/service"

	"github.com/google/uuid"
	"github.com/mitchellh/hashstructure/v2"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"trellis.tech/trellis/common.v1/clients/etcd"
	"trellis.tech/trellis/common.v1/crypto/base64"
	"trellis.tech/trellis/common.v1/errcode"
	"trellis.tech/trellis/common.v1/json"
	"trellis.tech/trellis/common.v1/logger"
)

const (
	defaultHeartbeat = time.Second * 5
	defaultTimeout   = defaultHeartbeat
)

type etcdRegistry struct {
	id      string
	options registry.Options

	Logger logger.Logger

	sync.RWMutex

	// map[registryFullPath]map[node.ID]id
	services map[string]uint64
	// map[registryFullPath]map[node.ID]leaseID
	leases map[string]clientv3.LeaseID
	//// map[registryFullPath]map[node.ID]worker
	workers map[string]*worker

	client etcd.Clientv3Facade
}

// NewRegistry new etcd registry
func NewRegistry(l logger.Logger, opts ...registry.Option) (registry.Registry, error) {

	p := &etcdRegistry{
		id: uuid.NewString(),

		Logger: l,

		services: make(map[string]uint64),
		leases:   make(map[string]clientv3.LeaseID),
		workers:  make(map[string]*worker),
	}

	for _, o := range opts {
		o(&p.options)
	}

	if err := configure(p); err != nil {
		return nil, err
	}

	return p, nil
}

// configure will setup the registry with new options
func configure(e *etcdRegistry) error {

	// setup the client
	cli, err := newClient(e)
	if err != nil {
		return err
	}

	if e.client != nil {
		e.client.Close()
	}

	// setup new client
	e.client = cli

	return nil
}

func (p *etcdRegistry) ID() string {
	return p.id
}

func (p *etcdRegistry) String() string {
	return registry.RegisterType_REGISTER_TYPE_ETCD.String()
}

func (p *etcdRegistry) Register(s *registry.ServiceNode) error {
	if s.GetService() == nil || s.GetService().GetName() == "" {
		return errcode.New("service name not found")
	}

	if s.GetNode() == nil || s.GetNode().GetValue() == "" {
		return errcode.New("service node's value not found")
	}

	fullRegPath := s.RegisteredServiceNode(p.options.Prefix)

	p.RLock()
	_, ok := p.workers[fullRegPath]
	p.RUnlock()
	if ok {
		return nil
	}

	heartbeat := time.Duration(s.GetNode().GetHeartbeat())
	if heartbeat <= 0 {
		heartbeat = defaultHeartbeat
	}

	ttl := time.Duration(s.GetNode().GetTTL())
	wer := &worker{
		node:        s,
		ticker:      time.NewTicker(heartbeat),
		fullRegPath: fullRegPath,
		stopSignal:  make(chan bool, 1),
		ttl:         ttl,
		timeout:     ttl,
	}

	if ttl <= 0 {
		wer.timeout = defaultTimeout
	}

	p.Logger.Debug("etd_register", "fullRegPath", fullRegPath, "Service", s)

	go func(wr *worker) {
		var count int
		for {
			fmt.Println(fullRegPath, *wr)
			if err := p.registerServiceNode(wr); err != nil {
				p.Logger.Warn("failed_and_retry_register", "worker", wr, "error", err.Error(),
					"retry_times", count, "max_retry_times", p.options.RetryTimes)
				if p.options.RetryTimes == 0 {
					continue
				}
				if p.options.RetryTimes <= count {
					panic(fmt.Errorf("%s regist into etcd failed times above: %d, %v", wr.fullRegPath, count, err))
				}
				count++
				continue
			}
			p.Logger.Debug("retry_register", "worker", wr,
				"retry_times", count, "max_retry_times", p.options.RetryTimes)

			count = 0
			select {
			case <-wr.stopSignal:
				return
			case <-wr.ticker.C:
				// nothing to do
				// TODO heartbeat
			}
		}
	}(wer)

	p.Lock()
	p.workers[fullRegPath] = wer
	p.Unlock()

	return nil
}

func (p *etcdRegistry) registerServiceNode(wr *worker) error {
	if wr == nil || wr.node == nil || wr.node.GetService().GetName() == "" ||
		wr.node.GetNode() == nil || wr.node.Node.Value == "" {
		return errcode.New("node should not be nil")
	}

	p.RLock()
	leaseID, ok := p.leases[wr.fullRegPath]
	p.RUnlock()

	p.Logger.Debug("register_service_node", "result", ok, "service", wr.node)

	if !ok {
		// minimum lease TTL is ttl-second
		ctx, cancel := context.WithTimeout(context.Background(), wr.timeout)
		defer cancel()
		resp, err := p.client.Get(ctx, wr.fullRegPath, clientv3.WithSerializable())
		if err != nil {
			return err
		}
		for _, kv := range resp.Kvs {
			if kv.Lease <= 0 {
				continue
			}
			leaseID = clientv3.LeaseID(kv.Lease)

			// decode the existing node
			srv := decode(kv.Value)
			if srv == nil || srv.Node == nil {
				continue
			}

			h, err := hashstructure.Hash(srv, hashstructure.FormatV2, nil)
			if err != nil {
				return err
			}

			// save the info
			p.Lock()
			p.leases[wr.fullRegPath] = leaseID
			p.services[wr.fullRegPath] = h
			p.Unlock()

			break
		}
	}

	var leaseNotFound bool
	if leaseID > 0 {
		if _, err := p.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
			if err != rpctypes.ErrLeaseNotFound {
				return err
			}
		}

		leaseNotFound = true
	}

	// create hash of service; uint64
	h, err := hashstructure.Hash(wr.node, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}

	// get existing hash for the service node
	p.RLock()
	v, ok := p.services[wr.fullRegPath]
	p.RUnlock()

	if ok && v == h && !leaseNotFound {
		return nil
	}

	var lgr *clientv3.LeaseGrantResponse
	if wr.ttl.Seconds() > 0 {

		ctx, cancel := context.WithTimeout(context.Background(), wr.ttl)
		defer cancel()

		// get a lease used to expire keys since we have a ttl
		lgr, err = p.client.Grant(ctx, int64(wr.ttl.Seconds()))
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), wr.timeout)
	defer cancel()
	var putOpts []clientv3.OpOption
	if lgr != nil {
		putOpts = append(putOpts, clientv3.WithLease(lgr.ID))
	}

	p.Logger.Debug("put_service_into_etcd", "path", wr.fullRegPath, "service", wr.node)

	if _, err = p.client.Put(ctx, wr.fullRegPath, encode(wr.node), putOpts...); err != nil {
		return err
	}

	p.Lock()
	// save our hash of the service
	p.services[wr.fullRegPath] = h
	// save our leaseID of the service
	if lgr != nil {
		p.leases[wr.fullRegPath] = lgr.ID
	}
	p.Unlock()

	return nil
}

//
func (p *etcdRegistry) Deregister(s *registry.ServiceNode) error {
	if s.Service.GetName() == "" {
		return errcode.New("service name not found")
	}

	fullRegPath := s.RegisteredServiceNode(p.options.Prefix)

	p.RLock()
	worker, ok := p.workers[fullRegPath]
	p.RUnlock()
	if !ok {
		return nil
	}

	p.Lock()
	defer p.Unlock()
	return p.stopWorker(worker)
}

func (p *etcdRegistry) Stop() error {
	p.Lock()
	defer p.Unlock()
	for _, w := range p.workers {
		if err := p.stopWorker(w); err != nil {
			return err
		}
	}

	if p.client != nil {
		p.client.Close()
	}

	return nil
}

func (p *etcdRegistry) Start() error {
	return nil
}

func (p *etcdRegistry) stopWorker(w *worker) error {

	w.stopSignal <- true
	close(w.stopSignal)

	delete(p.services, w.fullRegPath)
	delete(p.leases, w.fullRegPath)
	delete(p.workers, w.fullRegPath)

	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()
	_, err := p.client.Delete(ctx, w.fullRegPath)

	return err
}

func newClient(e *etcdRegistry, opts ...registry.Option) (etcd.Clientv3Facade, error) {
	return etcd.NewClient(*e.options.ETCDConfig)
}

func (p *etcdRegistry) Watch(s *service.Service) (registry.Watcher, error) {
	cli, err := newClient(p)
	if err != nil {
		return nil, err
	}
	return newEtcdWatcher(cli, p, s)
}

func encode(nn *registry.ServiceNode) string {
	bs, _ := json.Marshal(nn)
	return base64.Encode(base64.EncodeStd, bs)
}

func decode(bs []byte) *registry.ServiceNode {
	dst, err := base64.Decode(base64.EncodeStd, bs)
	if err != nil {
		return nil
	}

	var s *registry.ServiceNode
	if err = json.Unmarshal(dst, &s); err != nil {
		return nil
	}
	return s
}
