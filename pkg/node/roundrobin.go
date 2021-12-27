package node

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"trellis.tech/trellis/common.v1/types"
)

type roundrobin struct {
	Name string

	nodes   map[string]*Node
	indexes types.Strings

	count      int64
	robinIndex int64

	sync.RWMutex
}

// NewRoundRobin get roundrobin node manager
func NewRoundRobin(name string) (Manager, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, fmt.Errorf("name should not be nil")
	}
	return &roundrobin{Name: name}, nil
}

func (p *roundrobin) IsEmpty() bool {
	return atomic.LoadInt64(&p.count) == 0
}

func (p *roundrobin) Add(node *Node) {
	if node == nil {
		return
	}

	if node.Weight != 0 {
		node.Weight = 0
	}

	p.Lock()
	defer p.Unlock()
	p.add(node)
}

func (p *roundrobin) add(pNode *Node) {
	if p.nodes == nil {
		p.nodes = make(map[string]*Node)
	}

	_, ok := p.nodes[pNode.Value]
	if !ok {
		p.count++
		p.indexes = append(p.indexes, pNode.Value)

		sort.Sort(p.indexes)
	}

	p.nodes[pNode.Value] = pNode
}

func (p *roundrobin) NodeFor(...string) (*Node, bool) {
	if p.IsEmpty() {
		return nil, false
	}
	p.RLock()
	defer p.RUnlock()

	if p.robinIndex >= p.count {
		p.robinIndex = 0
	}
	node := p.nodes[p.indexes[int(p.robinIndex%p.count)]]

	p.robinIndex++

	return node, true
}

func (p *roundrobin) Remove() {
	p.Lock()
	defer p.Unlock()
	p.remove()
}

func (p *roundrobin) remove() {
	p.nodes = nil
	p.indexes = nil
	p.count = 0
	p.robinIndex = 0
}

func (p *roundrobin) RemoveByValue(value string) {
	p.Lock()
	defer p.Unlock()
	p.removeByValue(value)
}

func (p *roundrobin) removeByValue(id string) {
	if p.IsEmpty() {
		p.remove()
		return
	}

	_, ok := p.nodes[id]
	if !ok {
		return
	}
	delete(p.nodes, id)
	p.removeKey(id)
	p.count--

	sort.Sort(p.indexes)
}

func (p *roundrobin) removeKey(key string) {
	for i, v := range p.indexes {
		if v == key {
			p.indexes = append(p.indexes[:i], p.indexes[i+1:]...)
			break
		}
	}
}

func (p *roundrobin) PrintNodes() {
	p.RLock()
	defer p.RUnlock()

	for i, v := range p.nodes {
		fmt.Println("nodes:", i, *v)
	}
}
