package etcd

import (
	"time"

	"trellis.tech/trellis.v1/pkg/service"
)

type worker struct {
	service *service.Node

	fullServiceName string
	fullRegPath     string

	stopSignal chan bool

	// invoke self-register with ticker
	ticker *time.Ticker

	timeout time.Duration
	ttl     time.Duration
}
