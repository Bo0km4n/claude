// Code generated by protoc-gen-go. DO NOT EDIT.
// source: lr.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LRClient is the client API for LR service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LRClient interface {
	Heartbeat(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type lRClient struct {
	cc *grpc.ClientConn
}

func NewLRClient(cc *grpc.ClientConn) LRClient {
	return &lRClient{cc}
}

func (c *lRClient) Heartbeat(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.LR/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LR service

type LRServer interface {
	Heartbeat(context.Context, *Empty) (*Empty, error)
}

func RegisterLRServer(s *grpc.Server, srv LRServer) {
	s.RegisterService(&_LR_serviceDesc, srv)
}

func _LR_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LRServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LR/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LRServer).Heartbeat(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _LR_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.LR",
	HandlerType: (*LRServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Heartbeat",
			Handler:    _LR_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lr.proto",
}

func init() { proto.RegisterFile("lr.proto", fileDescriptor_lr_3661a2edbcd7620f) }

var fileDescriptor_lr_3661a2edbcd7620f = []byte{
	// 81 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xc8, 0x29, 0xd2, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0xdc, 0xa9, 0xb9, 0x05, 0x25, 0x95, 0x10,
	0x31, 0x23, 0x5d, 0x2e, 0x26, 0x9f, 0x20, 0x21, 0x75, 0x2e, 0x4e, 0x8f, 0xd4, 0xc4, 0xa2, 0x92,
	0xa4, 0xd4, 0xc4, 0x12, 0x21, 0x1e, 0x88, 0x94, 0x9e, 0x2b, 0x48, 0x99, 0x14, 0x0a, 0x2f, 0x89,
	0x0d, 0xcc, 0x31, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x5d, 0x84, 0x15, 0x10, 0x55, 0x00, 0x00,
	0x00,
}