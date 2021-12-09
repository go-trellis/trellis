package server

import (
	"context"

	"trellis.tech/trellis.v1/pkg/lifecycle"
	"trellis.tech/trellis.v1/pkg/message"
)

type Server interface {
	lifecycle.LifeCycle

	// Handle allows a single request to be made
	Handle(context.Context, *message.Request) (*message.Response, error)
	//// Stream is a bidirectional stream
	//Stream(Poster_StreamServer) error
	//// Publish publishes a payload and returns an empty payload
	//Publish(context.Context, *message.Payload) (*message.Payload, error)
}
