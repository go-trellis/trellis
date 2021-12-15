package registry

import (
	"time"

	"trellis.tech/trellis/common.v0/clients/etcd"
)

// Option initial options' functions
type Option func(*Options)

// Options new registry Options
type Options struct {
	Prefix     string
	ETCDConfig *etcd.Config
	Heartbeat  time.Duration
	TTL        time.Duration
	RetryTimes uint32
}

func Prefix(pre string) Option {
	return func(o *Options) {
		o.Prefix = pre
	}
}

func EtcdConfig(c *etcd.Config) Option {
	return func(o *Options) {
		o.ETCDConfig = c
	}
}

func Heartbeat(hb time.Duration) Option {
	return func(o *Options) {
		o.Heartbeat = hb
	}
}

func TTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.TTL = ttl
	}
}

func RetryTimes(times uint32) Option {
	return func(o *Options) {
		o.RetryTimes = times
	}
}
