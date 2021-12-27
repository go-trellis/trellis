package registry

import (
	"trellis.tech/trellis/common.v1/clients/etcd"
)

// Option initial options' functions
type Option func(*Options)

// Options new registry Options
type Options struct {
	Prefix     string
	ETCDConfig *etcd.Config
	RetryTimes int
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

func RetryTimes(times int) Option {
	return func(o *Options) {
		o.RetryTimes = times
	}
}
