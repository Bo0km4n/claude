// Code generated by protoc-gen-go. DO NOT EDIT.
// source: peer.proto

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

type NoticeFromLRRequest struct {
	GrpcPort             string   `protobuf:"bytes,1,opt,name=grpc_port,json=grpcPort" json:"grpc_port,omitempty"`
	TcpPort              string   `protobuf:"bytes,2,opt,name=tcp_port,json=tcpPort" json:"tcp_port,omitempty"`
	UdpPort              string   `protobuf:"bytes,3,opt,name=udp_port,json=udpPort" json:"udp_port,omitempty"`
	Addr                 string   `protobuf:"bytes,4,opt,name=addr" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NoticeFromLRRequest) Reset()         { *m = NoticeFromLRRequest{} }
func (m *NoticeFromLRRequest) String() string { return proto.CompactTextString(m) }
func (*NoticeFromLRRequest) ProtoMessage()    {}
func (*NoticeFromLRRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_peer_3acc5add53c04966, []int{0}
}
func (m *NoticeFromLRRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoticeFromLRRequest.Unmarshal(m, b)
}
func (m *NoticeFromLRRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoticeFromLRRequest.Marshal(b, m, deterministic)
}
func (dst *NoticeFromLRRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoticeFromLRRequest.Merge(dst, src)
}
func (m *NoticeFromLRRequest) XXX_Size() int {
	return xxx_messageInfo_NoticeFromLRRequest.Size(m)
}
func (m *NoticeFromLRRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NoticeFromLRRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NoticeFromLRRequest proto.InternalMessageInfo

func (m *NoticeFromLRRequest) GetGrpcPort() string {
	if m != nil {
		return m.GrpcPort
	}
	return ""
}

func (m *NoticeFromLRRequest) GetTcpPort() string {
	if m != nil {
		return m.TcpPort
	}
	return ""
}

func (m *NoticeFromLRRequest) GetUdpPort() string {
	if m != nil {
		return m.UdpPort
	}
	return ""
}

func (m *NoticeFromLRRequest) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func init() {
	proto.RegisterType((*NoticeFromLRRequest)(nil), "proto.NoticeFromLRRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PeerClient is the client API for Peer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PeerClient interface {
	NoticeFromLRRPC(ctx context.Context, in *NoticeFromLRRequest, opts ...grpc.CallOption) (*Empty, error)
}

type peerClient struct {
	cc *grpc.ClientConn
}

func NewPeerClient(cc *grpc.ClientConn) PeerClient {
	return &peerClient{cc}
}

func (c *peerClient) NoticeFromLRRPC(ctx context.Context, in *NoticeFromLRRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.Peer/NoticeFromLRRPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Peer service

type PeerServer interface {
	NoticeFromLRRPC(context.Context, *NoticeFromLRRequest) (*Empty, error)
}

func RegisterPeerServer(s *grpc.Server, srv PeerServer) {
	s.RegisterService(&_Peer_serviceDesc, srv)
}

func _Peer_NoticeFromLRRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoticeFromLRRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServer).NoticeFromLRRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Peer/NoticeFromLRRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServer).NoticeFromLRRPC(ctx, req.(*NoticeFromLRRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Peer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Peer",
	HandlerType: (*PeerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NoticeFromLRRPC",
			Handler:    _Peer_NoticeFromLRRPC_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "peer.proto",
}

func init() { proto.RegisterFile("peer.proto", fileDescriptor_peer_3acc5add53c04966) }

var fileDescriptor_peer_3acc5add53c04966 = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x48, 0x4d, 0x2d,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0xdc, 0xa9, 0xb9, 0x05, 0x25,
	0x95, 0x10, 0x31, 0xa5, 0x1a, 0x2e, 0x61, 0xbf, 0xfc, 0x92, 0xcc, 0xe4, 0x54, 0xb7, 0xa2, 0xfc,
	0x5c, 0x9f, 0xa0, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x69, 0x2e, 0xce, 0xf4, 0xa2,
	0x82, 0xe4, 0xf8, 0x82, 0xfc, 0xa2, 0x12, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x0e, 0x90,
	0x40, 0x40, 0x7e, 0x51, 0x89, 0x90, 0x24, 0x17, 0x47, 0x49, 0x72, 0x01, 0x44, 0x8e, 0x09, 0x2c,
	0xc7, 0x5e, 0x92, 0x5c, 0x00, 0x93, 0x2a, 0x4d, 0x81, 0x4a, 0x31, 0x43, 0xa4, 0x4a, 0x53, 0x20,
	0x52, 0x42, 0x5c, 0x2c, 0x89, 0x29, 0x29, 0x45, 0x12, 0x2c, 0x60, 0x61, 0x30, 0xdb, 0xc8, 0x99,
	0x8b, 0x25, 0x20, 0x35, 0xb5, 0x48, 0xc8, 0x9a, 0x8b, 0x1f, 0xc5, 0x15, 0x01, 0xce, 0x42, 0x52,
	0x10, 0x07, 0xea, 0x61, 0x71, 0x9d, 0x14, 0x0f, 0x54, 0xce, 0x15, 0xe4, 0x91, 0x24, 0x36, 0x30,
	0xc7, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xe9, 0x98, 0xa9, 0xa6, 0xeb, 0x00, 0x00, 0x00,
}
