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

type NoticeFromProxyRequest struct {
	Id                   uint32   `protobuf:"varint,5,opt,name=id" json:"id,omitempty"`
	GrpcPort             string   `protobuf:"bytes,1,opt,name=grpc_port,json=grpcPort" json:"grpc_port,omitempty"`
	TcpPort              string   `protobuf:"bytes,2,opt,name=tcp_port,json=tcpPort" json:"tcp_port,omitempty"`
	UdpPort              string   `protobuf:"bytes,3,opt,name=udp_port,json=udpPort" json:"udp_port,omitempty"`
	Addr                 string   `protobuf:"bytes,4,opt,name=addr" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NoticeFromProxyRequest) Reset()         { *m = NoticeFromProxyRequest{} }
func (m *NoticeFromProxyRequest) String() string { return proto.CompactTextString(m) }
func (*NoticeFromProxyRequest) ProtoMessage()    {}
func (*NoticeFromProxyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_peer_aaeca620418a948c, []int{0}
}
func (m *NoticeFromProxyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoticeFromProxyRequest.Unmarshal(m, b)
}
func (m *NoticeFromProxyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoticeFromProxyRequest.Marshal(b, m, deterministic)
}
func (dst *NoticeFromProxyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoticeFromProxyRequest.Merge(dst, src)
}
func (m *NoticeFromProxyRequest) XXX_Size() int {
	return xxx_messageInfo_NoticeFromProxyRequest.Size(m)
}
func (m *NoticeFromProxyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NoticeFromProxyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NoticeFromProxyRequest proto.InternalMessageInfo

func (m *NoticeFromProxyRequest) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *NoticeFromProxyRequest) GetGrpcPort() string {
	if m != nil {
		return m.GrpcPort
	}
	return ""
}

func (m *NoticeFromProxyRequest) GetTcpPort() string {
	if m != nil {
		return m.TcpPort
	}
	return ""
}

func (m *NoticeFromProxyRequest) GetUdpPort() string {
	if m != nil {
		return m.UdpPort
	}
	return ""
}

func (m *NoticeFromProxyRequest) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type GetPeerIDResponse struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPeerIDResponse) Reset()         { *m = GetPeerIDResponse{} }
func (m *GetPeerIDResponse) String() string { return proto.CompactTextString(m) }
func (*GetPeerIDResponse) ProtoMessage()    {}
func (*GetPeerIDResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_peer_aaeca620418a948c, []int{1}
}
func (m *GetPeerIDResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPeerIDResponse.Unmarshal(m, b)
}
func (m *GetPeerIDResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPeerIDResponse.Marshal(b, m, deterministic)
}
func (dst *GetPeerIDResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPeerIDResponse.Merge(dst, src)
}
func (m *GetPeerIDResponse) XXX_Size() int {
	return xxx_messageInfo_GetPeerIDResponse.Size(m)
}
func (m *GetPeerIDResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPeerIDResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetPeerIDResponse proto.InternalMessageInfo

func (m *GetPeerIDResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*NoticeFromProxyRequest)(nil), "proto.NoticeFromProxyRequest")
	proto.RegisterType((*GetPeerIDResponse)(nil), "proto.GetPeerIDResponse")
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
	NoticeFromProxyRPC(ctx context.Context, in *NoticeFromProxyRequest, opts ...grpc.CallOption) (*Empty, error)
	GetPeerID(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetPeerIDResponse, error)
}

type peerClient struct {
	cc *grpc.ClientConn
}

func NewPeerClient(cc *grpc.ClientConn) PeerClient {
	return &peerClient{cc}
}

func (c *peerClient) NoticeFromProxyRPC(ctx context.Context, in *NoticeFromProxyRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.Peer/NoticeFromProxyRPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerClient) GetPeerID(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetPeerIDResponse, error) {
	out := new(GetPeerIDResponse)
	err := c.cc.Invoke(ctx, "/proto.Peer/GetPeerID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Peer service

type PeerServer interface {
	NoticeFromProxyRPC(context.Context, *NoticeFromProxyRequest) (*Empty, error)
	GetPeerID(context.Context, *Empty) (*GetPeerIDResponse, error)
}

func RegisterPeerServer(s *grpc.Server, srv PeerServer) {
	s.RegisterService(&_Peer_serviceDesc, srv)
}

func _Peer_NoticeFromProxyRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoticeFromProxyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServer).NoticeFromProxyRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Peer/NoticeFromProxyRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServer).NoticeFromProxyRPC(ctx, req.(*NoticeFromProxyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Peer_GetPeerID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServer).GetPeerID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Peer/GetPeerID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServer).GetPeerID(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Peer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Peer",
	HandlerType: (*PeerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NoticeFromProxyRPC",
			Handler:    _Peer_NoticeFromProxyRPC_Handler,
		},
		{
			MethodName: "GetPeerID",
			Handler:    _Peer_GetPeerID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "peer.proto",
}

func init() { proto.RegisterFile("peer.proto", fileDescriptor_peer_aaeca620418a948c) }

var fileDescriptor_peer_aaeca620418a948c = []byte{
	// 235 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0xd0, 0x4b, 0x4a, 0xc4, 0x40,
	0x10, 0x06, 0x60, 0x12, 0x33, 0x3a, 0x29, 0x1f, 0x60, 0x2d, 0x24, 0x46, 0x84, 0x61, 0xdc, 0xcc,
	0x6a, 0x16, 0xce, 0x09, 0xc4, 0x17, 0x6e, 0x24, 0xe4, 0x02, 0xa2, 0xe9, 0x42, 0x7a, 0x31, 0x56,
	0x59, 0xa9, 0x80, 0xb3, 0xf1, 0x0a, 0x5e, 0x59, 0xd2, 0xdd, 0x0a, 0x3e, 0x56, 0x49, 0xff, 0x5f,
	0x43, 0xfd, 0xd5, 0x00, 0x42, 0xa4, 0x4b, 0x51, 0x36, 0xc6, 0x49, 0xf8, 0xd4, 0xbb, 0xb4, 0x16,
	0xdb, 0xc4, 0x6c, 0xfe, 0x91, 0xc1, 0xd1, 0x3d, 0x9b, 0xef, 0xe8, 0x46, 0x79, 0xdd, 0x28, 0xbf,
	0x6d, 0x5a, 0x7a, 0x1d, 0xa8, 0x37, 0x3c, 0x80, 0xdc, 0xbb, 0x6a, 0x32, 0xcb, 0x16, 0xfb, 0x6d,
	0xee, 0x1d, 0x9e, 0x40, 0xf9, 0xac, 0xd2, 0x3d, 0x08, 0xab, 0x55, 0xd9, 0x2c, 0x5b, 0x94, 0xed,
	0x74, 0x0c, 0x1a, 0x56, 0xc3, 0x63, 0x98, 0x5a, 0x27, 0xd1, 0xf2, 0x60, 0x3b, 0xd6, 0xc9, 0x17,
	0x0d, 0x2e, 0xd1, 0x56, 0xa4, 0xc1, 0x45, 0x42, 0x28, 0x1e, 0x9d, 0xd3, 0xaa, 0x08, 0x71, 0xf8,
	0x9f, 0x9f, 0xc1, 0xe1, 0x2d, 0x59, 0x43, 0xa4, 0x77, 0x57, 0x2d, 0xf5, 0xc2, 0x2f, 0x3d, 0xa5,
	0x2e, 0x71, 0x68, 0xee, 0xdd, 0xf9, 0x3b, 0x14, 0xe3, 0x0d, 0xbc, 0x00, 0xfc, 0xdd, 0xbe, 0xb9,
	0xc4, 0xd3, 0xb8, 0xdc, 0xf2, 0xff, 0xc5, 0xea, 0xbd, 0xc4, 0xd7, 0xe3, 0x3b, 0xe0, 0x0a, 0xca,
	0xef, 0x79, 0xf8, 0x83, 0xea, 0x2a, 0x9d, 0xfe, 0xf4, 0x79, 0xda, 0x0e, 0xb0, 0xfa, 0x0c, 0x00,
	0x00, 0xff, 0xff, 0x0b, 0xb2, 0xbe, 0xa6, 0x5f, 0x01, 0x00, 0x00,
}
