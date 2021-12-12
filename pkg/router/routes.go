package router

import (
	"sync"

	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/registry/memory"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis/common.v0/logger"
)

type routes struct {
	conf Config

	Logger logger.Logger

	Registry registry.Registry

	managerLocker sync.RWMutex

	nodeManagers map[string]node.Manager
}

func (p *routes) Start() (err error) {
	switch p.conf.RegistryConfig.RegisterType {
	case registry.RegisterType_etcd:
	case registry.RegisterType_memory:
		fallthrough
	default:
		p.Registry, err = memory.NewRegistry(
			//registry.Logger(p.conf.Logger.With("register", "memory")),
			registry.Prefix(p.conf.RegistryConfig.RegisterPrefix),
		)
	}

	for _, s := range p.conf.RegistryConfig.RegisterServiceNodes {
		if err = p.Registry.Register(s); err != nil {
			return
		}
	}

	for _, s := range p.conf.RegistryConfig.WatchServices {
		if err = p.Watch(s); err != nil {
			return
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
	if s == nil {
		return nil, false
	}
	servicePath := s.FullPath()
	p.managerLocker.RLock()
	manager, ok := p.nodeManagers[servicePath]
	if !ok {
		return nil, false
	}

	return manager.NodeFor(keys...)
}

func (p *routes) Register(s *service.ServiceNode) error {
	return p.Registry.Register(s)
}

func (p *routes) Deregister(s *service.ServiceNode) error {
	return p.Registry.Deregister(s)
}

func (p *routes) Watch(s *registry.WatchService) error {
	watcher, err := p.Registry.Watch(s.GetService())
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

			servicePath := r.ServiceNode.GetService().FullPath()
			p.managerLocker.RLocker()
			manager, ok := p.nodeManagers[servicePath]
			p.managerLocker.RUnlock()

			if !ok {
				manager, err = node.New(nodeType, servicePath)
				if err != nil {
					p.Logger.Errorf("failed_watch_service", "new_node_manager", servicePath, "error", err)
					continue
				}
			}

			switch r.Type {
			case service.EventType_create, service.EventType_update:
				p.Logger.Errorf("watch_service", "add_service_node", r.ServiceNode)
				manager.Add(r.ServiceNode.GetNode())
				p.managerLocker.Unlock()
			case service.EventType_delete:
				p.Logger.Errorf("watch_service", "delete_service_node", r.ServiceNode)
				manager.RemoveByValue(r.ServiceNode.GetNode().GetValue())
			}
		}
	}(watcher, s.NodeType)
	return nil
}
