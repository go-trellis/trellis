package memory

import (
	"sync"
	"time"

	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/service"

	"github.com/google/uuid"
	"trellis.tech/trellis/common.v0/errcode"
	"trellis.tech/trellis/common.v0/logger"
)

var (
	sendEventTime = 10 * time.Millisecond
)

type memory struct {
	id string

	sync.RWMutex

	prefix string

	// map[service_name]map[service_node_value]*registry.Service
	services map[string]map[string]*service.ServiceNode
	watchers map[string]*Watcher
}

// NewRegistry 生成新对象
func NewRegistry(l logger.Logger, opts ...registry.Option) (registry.Registry, error) {

	var options registry.Options
	for _, o := range opts {
		o(&options)
	}

	r := &memory{
		id: uuid.New().String(),

		prefix: options.Prefix,

		// domain/service version
		services: make(map[string]map[string]*service.ServiceNode),
		watchers: make(map[string]*Watcher),
	}

	return r, nil
}

// Start initial register
func (p *memory) Start() error {
	return nil
}

func (p *memory) Register(s *service.ServiceNode) error {
	p.Lock()
	defer p.Unlock()
	serviceName := s.GetService().GetPath(p.prefix)
	nodes, ok := p.services[serviceName]
	if !ok || nodes == nil {
		nodes = make(map[string]*service.ServiceNode)
	}

	nodeID := service.ReplaceURL(s.GetNode().GetValue())
	nodes[nodeID] = s

	p.services[serviceName] = nodes

	go p.sendEvent(&registry.Result{
		ID:          p.id,
		Timestamp:   time.Now(),
		Type:        service.EventType_update,
		ServiceNode: s})

	return nil
}

func (p *memory) Deregister(s *service.ServiceNode) error {
	p.Lock()
	defer p.Unlock()
	serviceName := s.GetService().GetPath(p.prefix)
	nodes, ok := p.services[serviceName]
	if !ok {
		return nil
	}

	nodeID := service.ReplaceURL(s.GetNode().GetValue())
	if _, ok := nodes[nodeID]; ok {
		delete(p.services[serviceName], nodeID)
	}

	if len(nodes) == 0 {
		delete(p.services, serviceName)
	}

	go p.sendEvent(&registry.Result{
		ID:          p.id,
		Timestamp:   time.Now(),
		Type:        service.EventType_delete,
		ServiceNode: s})

	return nil
}

func (p *memory) Watch(s *service.Service) (registry.Watcher, error) {

	if s == nil {
		return nil, errcode.New("watch unknown service (nil)")
	}

	w := &Watcher{
		id:     uuid.New().String(),
		exit:   make(chan bool),
		res:    make(chan *registry.Result),
		prefix: p.prefix,
		serv:   s,
	}

	p.Lock()
	p.watchers[w.id] = w
	p.Unlock()

	return w, nil
}

func (p *memory) Stop() error {
	return nil
}

func (p *memory) ID() string {
	return p.id
}

func (p *memory) String() string {
	return registry.RegisterType_memory.String()
}

func (p *memory) sendEvent(r *registry.Result) {
	p.RLock()
	watchers := make([]*Watcher, 0, len(p.watchers))
	for _, w := range p.watchers {
		watchers = append(watchers, w)
	}
	p.RUnlock()

	for _, w := range watchers {
		select {
		case <-w.exit:
			p.Lock()
			delete(p.watchers, w.id)
			p.Unlock()
		default:
			select {
			case w.res <- r:
			case <-time.After(sendEventTime):
			}
		}
	}
}
