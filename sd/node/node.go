package node

import (
	"fmt"

	"trellis.tech/trellis/common.v1/config"
)

// Type define node type
type Type uint8

// NodeType
const (
	NodeTypeDirect Type = iota
	NodeTypeRandom
	NodeTypeConsistent
	NodeTypeRoundRobin
)

// Node params for a node
type Node struct {
	// for recognize node with input id
	ID string `yaml:"id" json:"id"`
	// node's probability weight, roundrobin does not support
	Weight uint32 `yaml:"weight" json:"weight"`
	// node's value
	Value string `yaml:"value" json:"value"`
	// kvs for meta data
	Metadata config.Options `yaml:"options" json:"options"`

	number uint32
}

// Get get value from metadata
func (p *Node) Get(key string) (interface{}, bool) {
	if p.Metadata == nil {
		return nil, false
	}
	value, ok := p.Metadata[key]
	return value, ok
}

// Set set kv pair from metadata
func (p *Node) Set(key string, value interface{}) {
	if p.Metadata == nil {
		p.Metadata = config.Options{}
	}
	p.Metadata[key] = value
}

// Manager node manager functions defines.
type Manager interface {
	// adds a node to the node ring.
	Add(node *Node)
	// get the node responsible for the data key.
	NodeFor(keys ...string) (*Node, bool)
	// removes all nodes from the node ring.
	Remove()
	// removes a node from the node ring.
	RemoveByID(id string)
	// print all nodes
	PrintNodes()
	// is the node ring empty
	IsEmpty() bool
}

// New new node manager by node type, it has no nodes
func New(nt Type, name string) (Manager, error) {
	switch nt {
	case NodeTypeDirect:
		return NewDirect(name)
	case NodeTypeRandom:
		return NewRadmon(name)
	case NodeTypeConsistent:
		return NewConsistent(name)
	case NodeTypeRoundRobin:
		return NewRoundRobin(name)
	default:
		return nil, fmt.Errorf("not supperted type: %d", nt)
	}
}

// NewWithNodes new node manager by node type with nodes
func NewWithNodes(nt Type, name string, nodes []*Node) (Manager, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("nodes should at least one")
	}

	m, err := New(nt, name)
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		m.Add(n)
	}
	return m, nil
}
