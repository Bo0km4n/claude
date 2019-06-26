package service

import (
	"context"
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"google.golang.org/grpc"
)

// global variable
type remoteProxy struct {
	ID       uint32
	Addr     string
	GrpcPort string
	TcpPort  string
	UdpPort  string
}

var Protocol string
var NetConn net.Conn
var RemoteProxy remoteProxy
var IsCompletedJoinToProxy bool
var PeerSvc *PeerService

type PeerService struct {
	Seed string
	ID   string
}

func GetProxyTCPAddr() string {
	return RemoteProxy.Addr + ":" + RemoteProxy.TcpPort
}

func (p *PeerService) NoticeFromProxyRPC(ctx context.Context, in *proto.NoticeFromProxyRequest) (*proto.Empty, error) {
	if in.Addr == "" || in.GrpcPort == "" || in.TcpPort == "" || in.UdpPort == "" {
		return nil, errors.New("Proxy information is invalid")
	}

	if RemoteProxy.Addr != "" && RemoteProxy.GrpcPort != "" && RemoteProxy.TcpPort != "" && RemoteProxy.UdpPort != "" {
		return &proto.Empty{}, nil
	}

	RemoteProxy.ID = in.Id
	RemoteProxy.Addr = in.Addr
	RemoteProxy.GrpcPort = in.GrpcPort
	RemoteProxy.TcpPort = in.TcpPort
	RemoteProxy.UdpPort = in.UdpPort

	log.Printf("Registered Proxy | Addr: %s, GrpcPort: %s\n", RemoteProxy.Addr, RemoteProxy.GrpcPort)

	// set peer id
	id, err := getPeerIDString(&RemoteProxy, p.Seed)
	if err != nil {
		return nil, err
	}
	p.ID = id
	IsCompletedJoinToProxy = true

	return &proto.Empty{}, nil
}

func SetProxyInformation(seed string, iface string, proxyAddr string, useMulticast bool) {
	done := make(chan struct{})
	go LaunchGRPCService(done, seed)
	<-done

	time.Sleep(2)
	if useMulticast {
		UDPBcast(iface)
	} else {
		UDPUnicast(iface, proxyAddr)
	}
	for {
		if IsCompletedJoinToProxy {
			return
		}
		time.Sleep(1)
	}
}

func LaunchGRPCService(done chan struct{}, seed string) {
	port, err := net.Listen("tcp", ":"+"50051")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	PeerSvc = &PeerService{
		Seed: seed,
	}
	proto.RegisterPeerServer(
		server,
		PeerSvc,
	)

	log.Println("Start grpc services...")
	done <- struct{}{}
	server.Serve(port)
}

// Only ipv4
func extractPort(addr string) string {
	return strings.Split(addr, ":")[1]
}
