package node

import (
	"fmt"

	"trellis.tech/trellis/common.v1/config"
)

type Node struct {
	BaseNode `yaml:",inline" json:",inline"`
	Metadata config.Options `yaml:"metadata" json:"metadata"`
}

// Get value from metadata
func (p *Node) Get(key string) (interface{}, bool) {
	if p.Metadata == nil {
		return nil, false
	}
	value, ok := p.Metadata[key]
	return value, ok
}

// Copy node
func (p *Node) Copy() *Node {
	if p == nil {
		return nil
	}
	nd := &Node{
		Metadata: p.Metadata,
		BaseNode: BaseNode{},
	}
	return nd
}

// Set kv pair from metadata
func (p *Node) Set(key string, value interface{}) {
	if p.Metadata == nil {
		p.Metadata = make(map[string]interface{})
	}
	p.Metadata[key] = value
}

// Manager node manager functions defines.
type Manager interface {
	// Add adds a node to the node ring.
	Add(node *Node)
	// NodeFor get the node responsible for the data key.
	NodeFor(keys ...string) (*Node, bool)
	// Remove removes all nodes from the node ring.
	Remove()
	// RemoveByValue removes a node from the node ring.
	RemoveByValue(id string)
	// PrintNodes print all nodes
	PrintNodes()
	// IsEmpty is the node ring empty
	IsEmpty() bool
}

// New node manager with node type, it has no nodes
func New(nt NodeType, name string) (Manager, error) {
	switch nt {
	case NodeType_Direct:
		return NewDirect(name)
	case NodeType_Random:
		return NewRandom(name)
	case NodeType_Consistent:
		return NewConsistent(name)
	case NodeType_RoundRobin:
		return NewRoundRobin(name)
	default:
		return nil, fmt.Errorf("not supperted type: %d", nt)
	}
}

// NewWithNodes new node manager by node type with nodes
func NewWithNodes(nt NodeType, name string, nodes []*Node) (Manager, error) {
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
