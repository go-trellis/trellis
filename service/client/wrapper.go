package client

import (
	"context"

	"trellis.tech/go-trellis/common.v0/node"
)

// CallFunc represents the individual call func
type CallFunc func(ctx context.Context, node *node.Node, req Request, rsp interface{}, opts CallOptions) error

// Wrapper wraps a client and returns a client
type Wrapper func(Client) Client

// CallWrapper is a low level wrapper for the CallFunc
type CallWrapper func(CallFunc) CallFunc
