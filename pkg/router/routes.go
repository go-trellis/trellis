package router

import (
	"sync"

	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/registry/etcd"
	"trellis.tech/trellis.v1/pkg/registry/memory"
	"trellis.tech/trellis.v1/pkg/service"
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
	case registry.RegisterType_etcd:
		p.Registry, err = etcd.NewRegistry(
			p.Logger.With("registry", registry.RegisterType_etcd.String()),
			options...,
		)
	case registry.RegisterType_memory:
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
	p.managerLocker.RLock()
	manager, ok := p.nodeManagers[servicePath]
	p.managerLocker.RUnlock()
	if !ok {
		return nil, false
	}

	n, ok := manager.NodeFor(keys...)
	return n, ok
}

func (p *routes) Register(s *service.ServiceNode) error {
	return p.Registry.Register(s)
}

func (p *routes) Deregister(s *service.ServiceNode) error {
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
				r.ServiceNode.Node.Set("watch_service_config", s.Metadata)
			}

			switch r.Type {
			case service.EventType_create, service.EventType_update:
				p.Logger.Debug("watch_service", "add_service_node", r.ServiceNode)
				manager.Add(r.ServiceNode.Node)
			case service.EventType_delete:
				p.Logger.Debug("watch_service", "delete_service_node", r.ServiceNode)
				manager.RemoveByValue(r.ServiceNode.Node.GetValue())
			}
			p.managerLocker.RLock()
			p.nodeManagers[servicePath] = manager
			p.managerLocker.RUnlock()
		}
	}(watcher, s.NodeType)
	return nil
}
