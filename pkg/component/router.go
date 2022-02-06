package component

import (
	"sync"

	"trellis.tech/trellis.v1/pkg/service"

	"trellis.tech/trellis/common.v1/errcode"
)

var compR = &compRouter{
	newFuncs:   make(map[string]NewComponentFunc),
	components: make(map[string]Component),
}

var _ Router = (*compRouter)(nil)

type compRouter struct {
	mu sync.RWMutex
	// map[service]Component
	newFuncs map[string]NewComponentFunc

	components map[string]Component
}

func (p *compRouter) RegisterNewComponentFunc(s *service.Service, newFunc NewComponentFunc) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if newFunc == nil {
		delete(p.newFuncs, s.FullPath())
		return nil
	}
	if _, ok := p.newFuncs[s.FullPath()]; ok {
		return errcode.Newf("the new component function of service(%q) is already exist", s.FullPath())
	}
	p.newFuncs[s.FullPath()] = newFunc
	return nil
}

func (p *compRouter) RegisterComponent(s *service.Service, comp Component) error {
	if comp == nil {
		p.mu.Lock()
		delete(p.components, s.FullPath())
		p.mu.Unlock()
		return nil
	}
	p.mu.RLock()
	_, ok := p.components[s.FullPath()]
	p.mu.RUnlock()
	if ok {
		return errcode.Newf("service(%q) component is already exist", s.FullPath())
	}
	p.mu.Lock()
	p.components[s.FullPath()] = comp
	p.mu.Unlock()
	return comp.Start()
}

func (p *compRouter) NewComponent(c *Config) error {
	p.mu.RLock()
	newFunc, ok := p.newFuncs[c.Service.FullPath()]
	if !ok {
		p.mu.RUnlock()
		return errcode.Newf("new component function of service(%q) is not exist", c.Service.FullPath())
	}
	p.mu.RUnlock()

	comp, err := newFunc(c)
	if err != nil {
		return err
	}

	return p.RegisterComponent(c.Service, comp)
}

func (p *compRouter) GetComponent(s *service.Service) Component {
	p.mu.RLock()
	comp := p.components[s.FullPath()]
	p.mu.RUnlock()
	return comp
}

func (p *compRouter) StopComponents() error {
	p.mu.Lock()
	p.newFuncs = make(map[string]NewComponentFunc)
	components := p.components
	p.components = make(map[string]Component)
	p.mu.Unlock()

	var errs errcode.Errors
	wg := sync.WaitGroup{}
	for _, comp := range components {
		wg.Add(1)
		go func(comp Component) {
			defer wg.Done()
			err := comp.Stop()
			if err != nil {
				errs = append(errs, err)
			}
		}(comp)
	}
	wg.Wait()
	return errs.Errors()
}

func RegisterNewComponentFunc(s *service.Service, newFunc NewComponentFunc) {
	if err := compR.RegisterNewComponentFunc(s, newFunc); err != nil {
		panic(err)
	}
}

func RegisterComponent(s *service.Service, comp Component) {
	if err := compR.RegisterComponent(s, comp); err != nil {
		panic(err)
	}
}

func GetComponent(s *service.Service) Component {
	return compR.GetComponent(s)
}

func NewComponent(c *Config) error {
	return compR.NewComponent(c)
}

func StopComponents() error {
	return compR.StopComponents()
}
