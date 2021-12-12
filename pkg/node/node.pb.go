// Code generated by protoc-gen-go. DO NOT EDIT.
// source: trellis.tech/trellis.v1/proto/node.proto

package node

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
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

type NodeType int32

const (
	NodeType_Direct     NodeType = 0
	NodeType_Random     NodeType = 1
	NodeType_Consistent NodeType = 2
	NodeType_RoundRobin NodeType = 3
)

var NodeType_name = map[int32]string{
	0: "Direct",
	1: "Random",
	2: "Consistent",
	3: "RoundRobin",
}

var NodeType_value = map[string]int32{
	"Direct":     0,
	"Random":     1,
	"Consistent": 2,
	"RoundRobin": 3,
}

func (x NodeType) String() string {
	return proto.EnumName(NodeType_name, int32(x))
}

func (NodeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8475b92af2a52f7b, []int{0}
}

type Protocol int32

const (
	// 本地服务，直接访问
	Protocol_LOCAL Protocol = 0
	// RPC协议的服务
	Protocol_GRPC Protocol = 1
	// HTTP协议的服务
	Protocol_HTTP Protocol = 2
	// HTTP 3.0协议的服务
	Protocol_QUIC Protocol = 3
)

var Protocol_name = map[int32]string{
	0: "LOCAL",
	1: "GRPC",
	2: "HTTP",
	3: "QUIC",
}

var Protocol_value = map[string]int32{
	"LOCAL": 0,
	"GRPC":  1,
	"HTTP":  2,
	"QUIC":  3,
}

func (x Protocol) String() string {
	return proto.EnumName(Protocol_name, int32(x))
}

func (Protocol) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8475b92af2a52f7b, []int{1}
}

type Node struct {
	// @gotags: yaml:"weight"
	Weight uint32 `protobuf:"varint,1,opt,name=weight,proto3" json:"weight,omitempty" yaml:"weight"`
	// @gotags: yaml:"value"
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty" yaml:"value"`
	// @gotags: yaml:"ttl"
	TTL uint64 `protobuf:"varint,3,opt,name=TTL,proto3" json:"TTL,omitempty" yaml:"ttl"`
	// @gotags: yaml:"heartbeat"
	Heartbeat uint32 `protobuf:"varint,4,opt,name=heartbeat,proto3" json:"heartbeat,omitempty" yaml:"heartbeat"`
	// @gotags: yaml:"protocol"
	Protocol Protocol `protobuf:"varint,5,opt,name=protocol,proto3,enum=node.Protocol" json:"protocol,omitempty" yaml:"protocol"`
	// @gotags: yaml:"metadata,omitempty"
	Metadata map[string]string `protobuf:"bytes,6,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" yaml:"metadata,omitempty"`
	// @gotags: protobuf:"-" json:"-" yaml:"-"
	Number               uint32   `protobuf:"-" json:"-" yaml:"-"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_8475b92af2a52f7b, []int{0}
}

func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetWeight() uint32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *Node) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Node) GetTTL() uint64 {
	if m != nil {
		return m.TTL
	}
	return 0
}

func (m *Node) GetHeartbeat() uint32 {
	if m != nil {
		return m.Heartbeat
	}
	return 0
}

func (m *Node) GetProtocol() Protocol {
	if m != nil {
		return m.Protocol
	}
	return Protocol_LOCAL
}

func (m *Node) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Node) GetNumber() uint32 {
	if m != nil {
		return m.Number
	}
	return 0
}

func init() {
	proto.RegisterEnum("node.NodeType", NodeType_name, NodeType_value)
	proto.RegisterEnum("node.Protocol", Protocol_name, Protocol_value)
	proto.RegisterType((*Node)(nil), "node.Node")
	proto.RegisterMapType((map[string]string)(nil), "node.Node.MetadataEntry")
}

func init() {
	proto.RegisterFile("trellis.tech/trellis.v1/proto/node.proto", fileDescriptor_8475b92af2a52f7b)
}

var fileDescriptor_8475b92af2a52f7b = []byte{
	// 348 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x41, 0x6b, 0xea, 0x40,
	0x10, 0xc7, 0xdd, 0x24, 0xe6, 0x25, 0xf3, 0x50, 0x96, 0xe5, 0xf1, 0x58, 0x1e, 0xef, 0x10, 0x0a,
	0xa5, 0xc1, 0x43, 0xa4, 0xda, 0x43, 0xa9, 0xa7, 0x9a, 0x96, 0xb6, 0x60, 0xad, 0x5d, 0xd2, 0x4b,
	0x6f, 0x89, 0x19, 0x34, 0x18, 0xb3, 0x12, 0x57, 0x8b, 0x9f, 0xa5, 0x5f, 0xb6, 0x6c, 0xa2, 0x96,
	0x42, 0x7b, 0x59, 0x7e, 0xff, 0x99, 0xd9, 0x9d, 0xff, 0xec, 0x80, 0xaf, 0x4a, 0xcc, 0xf3, 0x6c,
	0x1d, 0x28, 0x9c, 0xce, 0xbb, 0x07, 0xb1, 0x3d, 0xef, 0xae, 0x4a, 0xa9, 0x64, 0xb7, 0x90, 0x29,
	0x06, 0x15, 0x32, 0x4b, 0xf3, 0xc9, 0xbb, 0x01, 0xd6, 0x58, 0xa6, 0xc8, 0xfe, 0x82, 0xfd, 0x86,
	0xd9, 0x6c, 0xae, 0x38, 0xf1, 0x88, 0xdf, 0x12, 0x7b, 0xc5, 0xfe, 0x40, 0x73, 0x1b, 0xe7, 0x1b,
	0xe4, 0x86, 0x47, 0x7c, 0x57, 0xd4, 0x82, 0x51, 0x30, 0xa3, 0x68, 0xc4, 0x4d, 0x8f, 0xf8, 0x96,
	0xd0, 0xc8, 0xfe, 0x83, 0x3b, 0xc7, 0xb8, 0x54, 0x09, 0xc6, 0x8a, 0x5b, 0xd5, 0x13, 0x9f, 0x01,
	0xd6, 0x01, 0xa7, 0xea, 0x3a, 0x95, 0x39, 0x6f, 0x7a, 0xc4, 0x6f, 0xf7, 0xda, 0x41, 0xe5, 0x65,
	0xb2, 0x8f, 0x8a, 0x63, 0x9e, 0x5d, 0x80, 0xb3, 0x44, 0x15, 0xa7, 0xb1, 0x8a, 0xb9, 0xed, 0x99,
	0xfe, 0xef, 0x1e, 0xaf, 0x6b, 0xb5, 0xcf, 0xe0, 0x71, 0x9f, 0xba, 0x2d, 0x54, 0xb9, 0x13, 0xc7,
	0x4a, 0xed, 0x7f, 0xbc, 0x59, 0x26, 0x58, 0xf2, 0x5f, 0xb5, 0xff, 0x5a, 0xfd, 0x1b, 0x40, 0xeb,
	0xcb, 0x15, 0x6d, 0x7d, 0x81, 0xbb, 0x6a, 0x4a, 0x57, 0x68, 0xfc, 0x7e, 0xc4, 0x2b, 0xe3, 0x92,
	0x74, 0x86, 0xe0, 0xe8, 0xa6, 0xd1, 0x6e, 0x85, 0x0c, 0xc0, 0xbe, 0xc9, 0x4a, 0x9c, 0x2a, 0xda,
	0xd0, 0x2c, 0xe2, 0x22, 0x95, 0x4b, 0x4a, 0x58, 0x1b, 0x20, 0x94, 0xc5, 0x3a, 0x5b, 0x2b, 0x2c,
	0x14, 0x35, 0xb4, 0x16, 0x72, 0x53, 0xa4, 0x42, 0x26, 0x59, 0x41, 0xcd, 0x4e, 0x1f, 0x9c, 0xc3,
	0x90, 0xcc, 0x85, 0xe6, 0xe8, 0x29, 0xbc, 0x1e, 0xd1, 0x06, 0x73, 0xc0, 0xba, 0x13, 0x93, 0x90,
	0x12, 0x4d, 0xf7, 0x51, 0x34, 0xa1, 0x86, 0xa6, 0xe7, 0x97, 0x87, 0x90, 0x9a, 0xc3, 0xb3, 0xd7,
	0xd3, 0x1f, 0x17, 0xb9, 0x98, 0x55, 0x6b, 0x1c, 0xe8, 0x23, 0xb1, 0xab, 0x6f, 0xeb, 0x7f, 0x04,
	0x00, 0x00, 0xff, 0xff, 0x90, 0x9a, 0x4e, 0x3a, 0xf8, 0x01, 0x00, 0x00,
}