// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: trellis.tech/trellis.v1/proto/message.proto

package message

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	service "trellis.tech/trellis.v1/pkg/service"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service *service.Service `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	// Request payload
	Payload *Payload `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trellis_tech_trellis_v1_proto_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_trellis_tech_trellis_v1_proto_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_trellis_tech_trellis_v1_proto_message_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetService() *service.Service {
	if x != nil {
		return x.Service
	}
	return nil
}

func (x *Request) GetPayload() *Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @gotags: json:"code"
	Code   int64  `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	ErrMsg string `protobuf:"bytes,2,opt,name=err_msg,json=errMsg,proto3" json:"err_msg,omitempty"`
	// Response payload @gotags: json:"payload"
	Payload *Payload `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trellis_tech_trellis_v1_proto_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_trellis_tech_trellis_v1_proto_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_trellis_tech_trellis_v1_proto_message_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Response) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *Response) GetPayload() *Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

type Payload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header map[string]string `protobuf:"bytes,1,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body   []byte            `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Payload) Reset() {
	*x = Payload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trellis_tech_trellis_v1_proto_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Payload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Payload) ProtoMessage() {}

func (x *Payload) ProtoReflect() protoreflect.Message {
	mi := &file_trellis_tech_trellis_v1_proto_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Payload.ProtoReflect.Descriptor instead.
func (*Payload) Descriptor() ([]byte, []int) {
	return file_trellis_tech_trellis_v1_proto_message_proto_rawDescGZIP(), []int{2}
}

func (x *Payload) GetHeader() map[string]string {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Payload) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_trellis_tech_trellis_v1_proto_message_proto protoreflect.FileDescriptor

var file_trellis_tech_trellis_v1_proto_message_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x74, 0x72, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x74,
	0x72, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x2b, 0x74, 0x72, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e,
	0x74, 0x65, 0x63, 0x68, 0x2f, 0x74, 0x72, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x61, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a,
	0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x63, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x5f, 0x6d, 0x73,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x12,
	0x2a, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x8e, 0x01, 0x0a, 0x07,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x34, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x1a, 0x39, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x2d, 0x5a, 0x2b,
	0x74, 0x72, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x74, 0x72, 0x65,
	0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x3b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_trellis_tech_trellis_v1_proto_message_proto_rawDescOnce sync.Once
	file_trellis_tech_trellis_v1_proto_message_proto_rawDescData = file_trellis_tech_trellis_v1_proto_message_proto_rawDesc
)

func file_trellis_tech_trellis_v1_proto_message_proto_rawDescGZIP() []byte {
	file_trellis_tech_trellis_v1_proto_message_proto_rawDescOnce.Do(func() {
		file_trellis_tech_trellis_v1_proto_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_trellis_tech_trellis_v1_proto_message_proto_rawDescData)
	})
	return file_trellis_tech_trellis_v1_proto_message_proto_rawDescData
}

var file_trellis_tech_trellis_v1_proto_message_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_trellis_tech_trellis_v1_proto_message_proto_goTypes = []interface{}{
	(*Request)(nil),         // 0: message.Request
	(*Response)(nil),        // 1: message.Response
	(*Payload)(nil),         // 2: message.Payload
	nil,                     // 3: message.Payload.HeaderEntry
	(*service.Service)(nil), // 4: service.Service
}
var file_trellis_tech_trellis_v1_proto_message_proto_depIdxs = []int32{
	4, // 0: message.Request.service:type_name -> service.Service
	2, // 1: message.Request.payload:type_name -> message.Payload
	2, // 2: message.Response.payload:type_name -> message.Payload
	3, // 3: message.Payload.header:type_name -> message.Payload.HeaderEntry
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_trellis_tech_trellis_v1_proto_message_proto_init() }
func file_trellis_tech_trellis_v1_proto_message_proto_init() {
	if File_trellis_tech_trellis_v1_proto_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trellis_tech_trellis_v1_proto_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trellis_tech_trellis_v1_proto_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trellis_tech_trellis_v1_proto_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Payload); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_trellis_tech_trellis_v1_proto_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_trellis_tech_trellis_v1_proto_message_proto_goTypes,
		DependencyIndexes: file_trellis_tech_trellis_v1_proto_message_proto_depIdxs,
		MessageInfos:      file_trellis_tech_trellis_v1_proto_message_proto_msgTypes,
	}.Build()
	File_trellis_tech_trellis_v1_proto_message_proto = out.File
	file_trellis_tech_trellis_v1_proto_message_proto_rawDesc = nil
	file_trellis_tech_trellis_v1_proto_message_proto_goTypes = nil
	file_trellis_tech_trellis_v1_proto_message_proto_depIdxs = nil
}
