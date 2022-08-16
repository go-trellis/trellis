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

package memory

import (
	"sync"
	"time"

	"trellis.tech/trellis.v1/pkg/registry"
	"trellis.tech/trellis.v1/pkg/service"

	"github.com/google/uuid"
	"trellis.tech/trellis/common.v1/errcode"
	"trellis.tech/trellis/common.v1/logger"
)

var (
	sendEventTime = 10 * time.Millisecond
)

type memory struct {
	id string

	sync.RWMutex

	prefix string

	// map[service_name]map[service_node_value]*registry.Service
	services map[string]map[string]*registry.ServiceNode
	watchers map[string]*Watcher
}

// NewRegistry 生成新对象
func NewRegistry(l logger.Logger, opts ...registry.Option) (registry.Registry, error) {

	var options registry.Options
	for _, o := range opts {
		o(&options)
	}

	r := &memory{
		id: uuid.NewString(),

		prefix: options.Prefix,

		// domain/service version
		services: make(map[string]map[string]*registry.ServiceNode),
		watchers: make(map[string]*Watcher),
	}

	return r, nil
}

// Start initial register
func (p *memory) Start() error {
	return nil
}

func (p *memory) Register(s *registry.ServiceNode) error {
	p.Lock()
	defer p.Unlock()
	serviceName := s.Service.GetPath(p.prefix)
	nodes, ok := p.services[serviceName]
	if !ok || nodes == nil {
		nodes = make(map[string]*registry.ServiceNode)
	}

	nodeID := service.ReplaceURL(s.Node.GetValue())
	nodes[nodeID] = s

	p.services[serviceName] = nodes

	go p.sendEvent(&registry.Result{
		Id:          p.id,
		Timestamp:   time.Now().UnixNano(),
		EventType:   service.EventType_EVENT_TYPE_UPDATE,
		ServiceNode: s})

	return nil
}

func (p *memory) Deregister(s *registry.ServiceNode) error {
	p.Lock()
	defer p.Unlock()
	serviceName := s.Service.GetPath(p.prefix)
	nodes, ok := p.services[serviceName]
	if !ok {
		return nil
	}

	nodeID := service.ReplaceURL(s.Node.GetValue())
	if _, ok := nodes[nodeID]; ok {
		delete(p.services[serviceName], nodeID)
	}

	if len(nodes) == 0 {
		delete(p.services, serviceName)
	}

	go p.sendEvent(&registry.Result{
		Id:          p.id,
		Timestamp:   time.Now().UnixNano(),
		EventType:   service.EventType_EVENT_TYPE_DELETE,
		ServiceNode: s})

	return nil
}

func (p *memory) Watch(s *service.Service) (registry.Watcher, error) {

	if s == nil {
		return nil, errcode.New("watch unknown service (nil)")
	}

	w := &Watcher{
		id:     uuid.NewString(),
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
	return registry.RegisterType_REGISTER_TYPE_MEMORY.String()
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
