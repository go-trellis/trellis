package server

import (
	"trellis.tech/trellis.v1/pkg/lifecycle"
)

type Server interface {
	lifecycle.LifeCycle

	TrellisServer
}
