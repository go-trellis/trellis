// Code generated by protoc-gen-go. DO NOT EDIT.
// source: trellis.tech/trellis.v1/proto/server.proto

package server

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
	message "trellis.tech/trellis.v1/pkg/message"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() {
	proto.RegisterFile("trellis.tech/trellis.v1/proto/server.proto", fileDescriptor_ba7fd9ea3c4ff36f)
}

var fileDescriptor_ba7fd9ea3c4ff36f = []byte{
	// 135 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2a, 0x29, 0x4a, 0xcd,
	0xc9, 0xc9, 0x2c, 0xd6, 0x2b, 0x49, 0x4d, 0xce, 0xd0, 0x87, 0x71, 0xca, 0x0c, 0xf5, 0x0b, 0x8a,
	0xf2, 0x4b, 0xf2, 0xf5, 0x8b, 0x53, 0x8b, 0xca, 0x52, 0x8b, 0xf4, 0xc0, 0x1c, 0x21, 0x96, 0xf4,
	0xa2, 0x82, 0x64, 0x29, 0x6d, 0xfc, 0x3a, 0x72, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x21, 0x5a,
	0x8c, 0x2c, 0xb8, 0xd8, 0x43, 0x20, 0x2a, 0x84, 0x74, 0xb9, 0x58, 0x9c, 0x13, 0x73, 0x72, 0x84,
	0x04, 0xf4, 0x60, 0x4a, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0xa4, 0x04, 0x91, 0x44, 0x8a,
	0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x95, 0x18, 0x9c, 0xb4, 0xa3, 0x34, 0x71, 0x5a, 0x94, 0x9d, 0x0e,
	0x75, 0x98, 0x35, 0x84, 0x4a, 0x62, 0x03, 0xdb, 0x66, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x3f,
	0x1c, 0xab, 0x11, 0xce, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TrellisClient is the client API for Trellis service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TrellisClient interface {
	// Call allows a single request to be made
	Call(ctx context.Context, in *message.Request, opts ...grpc.CallOption) (*message.Response, error)
}

type trellisClient struct {
	cc *grpc.ClientConn
}

func NewTrellisClient(cc *grpc.ClientConn) TrellisClient {
	return &trellisClient{cc}
}

func (c *trellisClient) Call(ctx context.Context, in *message.Request, opts ...grpc.CallOption) (*message.Response, error) {
	out := new(message.Response)
	err := c.cc.Invoke(ctx, "/grpc.Trellis/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrellisServer is the server API for Trellis service.
type TrellisServer interface {
	// Call allows a single request to be made
	Call(context.Context, *message.Request) (*message.Response, error)
}

// UnimplementedTrellisServer can be embedded to have forward compatible implementations.
type UnimplementedTrellisServer struct {
}

func (*UnimplementedTrellisServer) Call(ctx context.Context, req *message.Request) (*message.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}

func RegisterTrellisServer(s *grpc.Server, srv TrellisServer) {
	s.RegisterService(&_Trellis_serviceDesc, srv)
}

func _Trellis_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(message.Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrellisServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Trellis/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrellisServer).Call(ctx, req.(*message.Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Trellis_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Trellis",
	HandlerType: (*TrellisServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _Trellis_Call_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "trellis.tech/trellis.v1/proto/server.proto",
}