package etcd

import (
	"time"

	"trellis.tech/trellis.v1/pkg/service"
)

type worker struct {
	service *service.ServiceNode

	fullServiceName string
	fullRegPath     string

	stopSignal chan bool

	// invoke self-register with ticker
	ticker *time.Ticker

	interval time.Duration
}
