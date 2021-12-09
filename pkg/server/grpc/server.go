package grpc

import (
	"context"

	"google.golang.org/grpc"

	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/server"
)

var _ server.TrellisServer = (*Server)(nil)

type Server struct {
}

func (p *Server) Call(context.Context, *message.Request) (*message.Response, error) {
	return nil, nil
}

func (p *Server) Stream(server.Trellis_StreamServer) error {
	return nil
}

func (p *Server) Publish(context.Context, *message.Payload) (*message.Payload, error) {
	return nil, nil
}

func NewServer() *Server {
	s := &Server{}
	rpcServer := grpc.NewServer()
	server.RegisterTrellisServer(rpcServer, s)
	return s
}
