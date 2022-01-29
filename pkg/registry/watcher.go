package registry

import (
	"time"

	"trellis.tech/trellis.v1/pkg/clients"
	"trellis.tech/trellis.v1/pkg/node"
	"trellis.tech/trellis.v1/pkg/service"
)

// Watcher is an interface that returns updates
// about services within the registry.
type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}

// Result is registry result
type Result struct {
	// ID is registry id
	ID string
	// Type defines type of event
	Type service.EventType
	// Timestamp is event timestamp
	Timestamp time.Time
	// ServiceNode is registered service node
	ServiceNode *service.ServiceNode
}

type WatchService struct {
	Service  *service.Service `yaml:"service" json:"service"`
	NodeType node.NodeType    `yaml:"node_type" json:"node_type"`

	Metadata *WatchServiceMetadata `yaml:"metadata" json:"metadata"`
}

type WatchServiceMetadata struct {
	ClientConfig *clients.Config `yaml:"client_config" json:"client_config"`
}
