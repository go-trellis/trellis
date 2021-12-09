/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

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

//
//import (
//	"context"
//	"fmt"
//	"sync"
//	"time"
//
//	"trellis.tech/trellis.v1/pkg/registry"
//	"trellis.tech/trellis.v1/pkg/service"
//	"trellis.tech/trellis/common.v0/clients/etcd"
//	"trellis.tech/trellis/common.v0/crypto/base64"
//	"trellis.tech/trellis/common.v0/errcode"
//	"trellis.tech/trellis/common.v0/json"
//	"trellis.tech/trellis/common.v0/node"
//
//	"github.com/google/uuid"
//	"github.com/mitchellh/hashstructure/v2"
//	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
//	clientv3 "go.etcd.io/etcd/client/v3"
//)
//
//type etcdRegistry struct {
//	id string
//	//options registry.Options
//
//	config etcd.Config
//
//	sync.RWMutex
//
//	// map[registryFullPath]map[node.ID]id
//	services map[string]uint64
//	// map[registryFullPath]map[node.ID]leaseID
//	leases map[string]clientv3.LeaseID
//	// map[registryFullPath]map[node.ID]worker
//	workers map[string]*worker
//
//	client etcd.Clientv3Facade
//}
//
//// type register map[string]uint64
//// type leases map[string]clientv3.LeaseID
//// type workers map[string]*worker
//
//// NewRegistry new etcd registry
//func NewRegistry(cfg etcd.Config) (registry.Registry, error) {
//
//	p := &etcdRegistry{
//		id: uuid.New().String(),
//
//		config: cfg,
//
//		services: make(map[string]uint64),
//		leases:   make(map[string]clientv3.LeaseID),
//		workers:  make(map[string]*worker),
//	}
//
//	if err := configure(p); err != nil {
//		return nil, err
//	}
//
//	return p, nil
//}
//
//// configure will setup the registry with new options
//func configure(e *etcdRegistry) error {
//
//	// setup the client
//	cli, err := etcd.NewClient(e.config)
//	if err != nil {
//		return err
//	}
//
//	if e.client != nil {
//		e.client.Close()
//	}
//
//	// setup new client
//	e.client = cli
//
//	return nil
//}
//
//func (p *etcdRegistry) ID() string {
//	return p.id
//}
//
//func (p *etcdRegistry) Register(s *service.Service) error {
//	if s.GetName() == "" {
//		return errcode.New("service name not found")
//	}
//
//	fullRegPath := s.FullRegistryPath()
//
//	p.RLock()
//	_, ok := p.workers[fullRegPath]
//	p.RUnlock()
//	if ok {
//		return nil
//	}
//
//	wer := &worker{
//		service: &registry.Service{
//			Service: *s,
//			Node: &node.Node{
//				ID:     s.ID(p.options.ServerAddr),
//				Value:  p.options.ServerAddr,
//				Weight: options.Weight,
//			},
//		},
//		ticker:      time.NewTicker(options.Heartbeat),
//		fullRegPath: fullRegPath,
//		stopSignal:  make(chan bool, 1),
//		options:     options,
//	}
//
//	p.options.Logger.Debug("etd_register", "fullRegPath", fullRegPath, "Service", s)
//
//	go func(wr *worker) {
//		var count uint32
//		for {
//			if err := p.registerServiceNode(wr); err != nil {
//				p.options.Logger.Warn("failed_and_retry_regsiter", "worker", wr, "error", err.Error(),
//					"retry_times", count, "max_retry_times", p.options.RetryTimes)
//				if p.options.RetryTimes == 0 {
//					continue
//				}
//				if p.options.RetryTimes <= count {
//					panic(fmt.Errorf("%s regist into etcd failed times above: %d, %v", wr.fullRegPath, count, err))
//				}
//				count++
//				continue
//			}
//			p.options.Logger.Debug("retry_regsiter", "worker", wr,
//				"retry_times", count, "max_retry_times", p.options.RetryTimes)
//
//			count = 0
//			select {
//			case <-wr.stopSignal:
//				return
//			case <-wr.ticker.C:
//				// nothing to do
//			}
//		}
//	}(wer)
//
//	p.Lock()
//	p.workers[fullRegPath] = wer
//	p.Unlock()
//
//	return nil
//}
//
//func (p *etcdRegistry) registerServiceNode(wr *worker) error {
//	if wr == nil || wr.service == nil || wr.service.GetName() == "" ||
//		wr.service.Node == nil || wr.service.Node.ID == "" {
//		return errcode.New("node should not be nil")
//	}
//
//	p.RLock()
//	leaseID, ok := p.leases[wr.fullRegPath]
//	p.RUnlock()
//
//	p.options.Logger.Debug("register_service_node", "result", ok, "service", wr.service)
//
//	if !ok {
//		// minimum lease TTL is ttl-second
//		ctx, cancel := context.WithTimeout(context.Background(), p.options.Timeout)
//		defer cancel()
//		resp, err := p.client.Get(ctx, wr.fullRegPath, clientv3.WithSerializable())
//		if err != nil {
//			return err
//		}
//		for _, kv := range resp.Kvs {
//			if kv.Lease <= 0 {
//				continue
//			}
//			leaseID = clientv3.LeaseID(kv.Lease)
//
//			// decode the existing node
//			srv := decode(kv.Value)
//			if srv == nil || srv.Node == nil {
//				continue
//			}
//
//			h, err := hashstructure.Hash(srv, hashstructure.FormatV2, nil)
//			if err != nil {
//				return err
//			}
//
//			// save the info
//			p.Lock()
//			p.leases[wr.fullRegPath] = leaseID
//			p.services[wr.fullRegPath] = h
//			p.Unlock()
//
//			break
//		}
//	}
//
//	var leaseNotFound bool
//	if leaseID > 0 {
//		if _, err := p.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
//			if err != rpctypes.ErrLeaseNotFound {
//				return err
//			}
//		}
//
//		leaseNotFound = true
//	}
//
//	// create hash of service; uint64
//	h, err := hashstructure.Hash(wr.service, hashstructure.FormatV2, nil)
//	if err != nil {
//		return err
//	}
//
//	// get existing hash for the service node
//	p.RLock()
//	v, ok := p.services[wr.fullRegPath]
//	p.RUnlock()
//
//	if ok && v == h && !leaseNotFound {
//		return nil
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), p.options.Timeout)
//	defer cancel()
//
//	var lgr *clientv3.LeaseGrantResponse
//	if wr.options.TTL.Seconds() > 0 {
//		// get a lease used to expire keys since we have a ttl
//		lgr, err = p.client.Grant(ctx, int64(wr.options.TTL.Seconds()))
//		if err != nil {
//			return err
//		}
//	}
//	var putOpts []clientv3.OpOption
//	if lgr != nil {
//		putOpts = append(putOpts, clientv3.WithLease(lgr.ID))
//	}
//
//	p.options.Logger.Debug("put_service_into_etcd", "path", wr.fullRegPath, "service", wr.service)
//
//	if _, err = p.client.Put(ctx, wr.fullRegPath, encode(wr.service), putOpts...); err != nil {
//		return err
//	}
//
//	p.Lock()
//	// save our hash of the service
//	p.services[wr.fullRegPath] = h
//	// save our leaseID of the service
//	if lgr != nil {
//		p.leases[wr.fullRegPath] = lgr.ID
//	}
//	p.Unlock()
//
//	return nil
//}
//
//func (p *etcdRegistry) Deregister(s *service.Service, opts ...registry.DeregisterOption) error {
//	if s.GetName() == "" {
//		return errcode.New("service name not found")
//	}
//
//	fullRegPath := s.FullRegistryPath(p.options.ServerAddr)
//
//	p.RLock()
//	worker, ok := p.workers[fullRegPath]
//	p.RUnlock()
//	if !ok {
//		return nil
//	}
//
//	p.Lock()
//	defer p.Unlock()
//	return p.stopWorker(worker)
//}
//
//func (p *etcdRegistry) String() string {
//	return service.RegisterType_name[int32(service.RegisterType_etcd)]
//}
//
//func (p *etcdRegistry) Stop() error {
//	p.Lock()
//	defer p.Unlock()
//	for _, w := range p.workers {
//		if err := p.stopWorker(w); err != nil {
//			return err
//		}
//	}
//
//	if p.client != nil {
//		p.client.Close()
//	}
//
//	return nil
//}
//
//func (p *etcdRegistry) stopWorker(w *worker) error {
//
//	w.stopSignal <- true
//	close(w.stopSignal)
//
//	delete(p.services, w.fullRegPath)
//	delete(p.leases, w.fullRegPath)
//	delete(p.workers, w.fullRegPath)
//
//	ctx, cancel := context.WithTimeout(context.Background(), p.options.Timeout)
//	defer cancel()
//	_, err := p.client.Delete(ctx, w.fullRegPath)
//
//	return err
//}
//
//func (p *etcdRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
//	cli, err := newClient(p)
//	if err != nil {
//		return nil, err
//	}
//	return newEtcdWatcher(cli, p.id, p.options.Timeout, opts...)
//}
//
//func encode(nn *registry.Service) string {
//	bs, _ := json.Marshal(nn)
//	return base64.Encode(base64.EncodeStd, bs)
//}
//
//func decode(bs []byte) *registry.Service {
//	dst, err := base64.Decode(base64.EncodeStd, bs)
//	if err != nil {
//		return nil
//	}
//
//	var s *registry.Service
//	json.Unmarshal(dst, &s)
//	return s
//}
