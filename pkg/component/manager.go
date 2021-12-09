package component

import (
	"sync"

	"trellis.tech/trellis/common.v0/errcode"

	"trellis.tech/trellis.v1/pkg/service"
)

var m = &manager{
	newFuncs:   make(map[string]NewComponentFunc),
	components: make(map[string]Component),
}

var _ Manager = (*manager)(nil)

type Manager interface {
	RegisterNewComponentFunc(s *service.Service, newFunc NewComponentFunc) error
	RegisterComponent(s *service.Service, component Component) error
	NewComponent(s *service.Service, opts ...Option) error
	GetComponent(*service.Service) Component
	Stop() error
}

func GetManager() Manager {
	return m
}

type manager struct {
	mu sync.RWMutex
	// map[service]Component
	newFuncs map[string]NewComponentFunc

	components map[string]Component
}

func (p *manager) RegisterNewComponentFunc(s *service.Service, newFunc NewComponentFunc) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if newFunc == nil {
		delete(p.newFuncs, s.FullPath())
		return nil
	}
	if _, ok := m.newFuncs[s.FullPath()]; ok {
		return errcode.Newf("the new component function of service(%q) is already exist", s.FullPath())
	}
	m.newFuncs[s.FullPath()] = newFunc
	return nil
}

func (p *manager) RegisterComponent(s *service.Service, component Component) error {
	if component == nil {
		p.mu.Lock()
		delete(p.components, s.FullPath())
		p.mu.Unlock()
		return nil
	}
	_, ok := p.components[s.FullPath()]
	if ok {
		return errcode.Newf("service(%q) component is already exist", s.FullPath())
	}
	p.mu.Lock()
	p.components[s.FullPath()] = component
	p.mu.Unlock()

	return component.Start()
}

func (p *manager) NewComponent(s *service.Service, opts ...Option) error {
	p.mu.RLock()
	newFunc, ok := m.newFuncs[s.FullPath()]
	if !ok {
		p.mu.RUnlock()
		return errcode.Newf("new component function of service(%q) is not exist", s.FullPath())
	}
	p.mu.RUnlock()

	component, err := newFunc(opts...)
	if err != nil {
		return err
	}

	p.mu.Lock()
	p.components[s.FullPath()] = component
	p.mu.Unlock()

	return component.Start()
}

func (p *manager) GetComponent(s *service.Service) Component {
	p.mu.RLock()
	component := m.components[s.FullPath()]
	p.mu.RUnlock()
	return component
}

func (p *manager) Stop() error {
	p.mu.Lock()
	p.newFuncs = make(map[string]NewComponentFunc)
	components := p.components
	p.components = make(map[string]Component)
	p.mu.Unlock()

	var errs errcode.Errors
	wg := sync.WaitGroup{}
	for _, component := range components {
		wg.Add(1)
		go func(component Component) {
			defer wg.Done()
			err := component.Stop()
			if err != nil {
				errs = append(errs, err)
			}
		}(component)
	}
	wg.Wait()
	return errs.Errors()
}
