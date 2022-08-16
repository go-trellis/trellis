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

package router

import (
	"context"
	"fmt"
	"sync"

	"trellis.tech/trellis.v1/pkg/clients/client"
	"trellis.tech/trellis.v1/pkg/clients/local"
	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/registry/etcd"
	"trellis.tech/trellis.v1/pkg/registry/memory"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/json"
	"trellis.tech/trellis/common.v1/logger"
)

type routes struct {
	conf Config

	Logger logger.Logger

	Registry registry.Registry

	managerLocker sync.RWMutex

	nodeManagers map[string]node.Manager
}

func (p *routes) Start() (err error) {
	options := []registry.Option{
		registry.Prefix(p.conf.RegistryConfig.RegisterPrefix),
		registry.EtcdConfig(&p.conf.RegistryConfig.ETCDConfig),
		registry.RetryTimes(p.conf.RegistryConfig.RetryTimes),
	}
	switch p.conf.RegistryConfig.RegisterType {
	case registry.RegisterType_REGISTER_TYPE_ETCD:
		p.Registry, err = etcd.NewRegistry(
			p.Logger.With("registry", registry.RegisterType_REGISTER_TYPE_ETCD.String()),
			options...,
		)
	case registry.RegisterType_REGISTER_TYPE_MEMORY:
		fallthrough
	default:
		p.Registry, err = memory.NewRegistry(
			p.Logger.With("registry", p.conf.RegistryConfig.RegisterType.String()),
			options...,
		)
	}

	for _, s := range p.conf.RegistryConfig.RegisterServiceNodes {
		if err = p.Registry.Register(s); err != nil {
			return
		}
	}

	for _, s := range p.conf.RegistryConfig.WatchServices {
		if err = p.Watch(s); err != nil {
			return err
		}
	}
	return
}

func (p *routes) Stop() error {
	if err := p.Registry.Stop(); err != nil {
		return err
	}
	return nil
}

func (p *routes) GetServiceNode(s *service.Service, keys ...string) (*node.Node, bool) {
	servicePath := s.FullPath()
	if servicePath == "" {
		return nil, false
	}
	p.managerLocker.RLock()
	manager, ok := p.nodeManagers[servicePath]
	fmt.Println("GetServiceNode", servicePath, ok)
	p.managerLocker.RUnlock()
	if !ok {
		return nil, false
	}

	n, ok := manager.NodeFor(keys...)
	return n, ok
}

func (p *routes) Register(s *registry.ServiceNode) error {
	return p.Registry.Register(s)
}

func (p *routes) Deregister(s *registry.ServiceNode) error {
	return p.Registry.Deregister(s)
}

func (p *routes) Watch(s *registry.WatchService) error {
	watcher, err := p.Registry.Watch(s.Service)
	if err != nil {
		return err
	}
	go func(wch registry.Watcher, nodeType node.NodeType) {
		for {
			r, err := wch.Next()
			if err != nil {
				p.Logger.Errorf("failed_watch_service", "error", err)
				continue
			}

			servicePath := r.ServiceNode.Service.FullPath()
			p.managerLocker.RLock()
			manager, ok := p.nodeManagers[servicePath]
			p.managerLocker.RUnlock()

			if !ok {
				manager, err = node.New(nodeType, servicePath)
				if err != nil {
					p.Logger.Errorf("failed_watch_service", "new_node_manager", servicePath, "error", err)
					continue
				}
			}

			if s.Metadata != nil {
				bs, _ := json.Marshal(s.Metadata)
				r.ServiceNode.Node.Set("watch_service_config", string(bs))
			}

			switch r.GetEventType() {
			case service.EventType_EVENT_TYPE_CREATE, service.EventType_EVENT_TYPE_UPDATE:
				p.Logger.Debug("watch_service", "add_service_node", r.ServiceNode)
				manager.Add(r.ServiceNode.Node)
			case service.EventType_EVENT_TYPE_DELETE:
				p.Logger.Debug("watch_service", "delete_service_node", r.ServiceNode)
				manager.RemoveByValue(r.ServiceNode.Node.Value)
			}
			p.managerLocker.RLock()
			p.nodeManagers[servicePath] = manager
			p.managerLocker.RUnlock()
		}
	}(watcher, s.NodeType)
	return nil
}

func (p *routes) GetCaller(s *service.Service, keys ...string) (server.Caller, []server.CallOption, error) {
	var (
		c    server.Caller
		opts []server.CallOption
		err  error
	)
	serviceNode, ok := p.GetServiceNode(s, s.String())
	if !ok {
		c, opts, err = local.NewClient()
	} else {
		c, opts, err = client.New(serviceNode)
	}
	if err != nil {
		return nil, nil, err
	}
	return c, opts, nil
}

func (p *routes) Call(ctx context.Context, msg *message.Request) (*message.Response, error) {
	var (
		c    server.Caller
		opts []server.CallOption
		err  error
	)

	fmt.Println("ss.......", msg.GetService())
	serviceNode, ok := p.GetServiceNode(msg.GetService(), msg.String())
	fmt.Println("serviceNode.......", serviceNode)
	if !ok {
		c, opts, err = local.NewClient()
	} else {
		c, opts, err = client.New(serviceNode)
	}
	if err != nil {
		return nil, err
	}
	return c.Call(ctx, msg, opts...)
}

//
//func (p *routes) Call(ctx context.Context, msg *message.Request) (*message.Response, error) {
//	var (
//		c    server.Caller
//		opts []server.CallOption
//		err  error
//	)
//	serviceNode, ok := p.GetServiceNode(msg.GetService(), msg.String())
//	if !ok {
//		c, opts, err = local.NewClient()
//	} else {
//		c, opts, err = client.New(serviceNode)
//	}
//	if err != nil {
//		return nil, err
//	}
//	return c.Call(ctx, msg, opts...)
//}
