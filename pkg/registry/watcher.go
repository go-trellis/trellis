package registry

import (
	"time"

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
