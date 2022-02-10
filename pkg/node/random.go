package node

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type random struct {
	Name string

	nodes map[string]*Node
	rings map[int64]*Node

	count int64

	sync.RWMutex
}

// NewRandom get random node manager
func NewRandom(name string) (Manager, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, fmt.Errorf("name should not be nil")
	}
	return &random{Name: name}, nil
}

func (p *random) IsEmpty() bool {
	return atomic.LoadInt64(&p.count) == 0
}

func (p *random) Add(node *Node) {
	if node == nil {
		return
	}
	p.Lock()
	defer p.Unlock()
	p.add(node)
}

func (p *random) add(pNode *Node) {
	if p.nodes == nil {
		p.nodes = make(map[string]*Node)
	}

	p.nodes[pNode.Value] = pNode

	p.updateRings()
}

func (p *random) Remove() {
	p.Lock()
	defer p.Unlock()
	p.remove()
}

func (p *random) remove() {
	p.nodes = nil
	p.rings = nil
	p.count = 0
}

func (p *random) RemoveByValue(id string) {
	p.Lock()
	defer p.Unlock()
	p.removeByValue(id)
}

func (p *random) removeByValue(id string) {
	if p.nodes == nil {
		return
	} else if p.IsEmpty() {
		p.remove()
		return
	}

	node := p.nodes[id]
	if node == nil {
		return
	}

	delete(p.nodes, id)
	p.updateRings()
}

func (p *random) NodeFor(...string) (*Node, bool) {
	p.RLock()
	defer p.RUnlock()
	if p.IsEmpty() {
		return nil, false
	}

	rand.Seed(time.Now().UnixNano())

	return p.rings[rand.Int63n(p.count)], true
}

func (p *random) updateRings() {
	p.rings = make(map[int64]*Node)

	p.count = 0
	for _, v := range p.nodes {

		for i := 0; i < int(v.Weight); i++ {
			ring := v.Copy()
			ring.Number = uint32(i + 1)
			p.rings[p.count] = ring
			p.count++
		}
	}
}

func (p *random) PrintNodes() {
	p.Lock()
	defer p.Unlock()

	for i, v := range p.nodes {
		fmt.Println("nodes:", i, v.Copy())
	}
}
