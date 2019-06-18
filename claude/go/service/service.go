package service

import (
	"context"
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"github.com/k0kubun/pp"
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

	// if err := peerJoin(); err != nil {
	// 	log.Fatalf("Failed join to Proxy: %+v\n", err)
	// }

	// set peer id
	id, err := getPeerIDString(&RemoteProxy, p.Seed)
	if err != nil {
		return nil, err
	}
	pp.Println(id, p.Seed)
	p.ID = id
	IsCompletedJoinToProxy = true

	return &proto.Empty{}, nil
}

func SetProxyInformation(seed string) {
	done := make(chan struct{})
	go LaunchGRPCService(done, seed)
	<-done

	time.Sleep(2)
	UDPBcast()
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

// func peerJoin() error {
// 	latitude, longitude := geo.GetLocation()
// 	request := &proto.PeerJoinRequest{
// 		PeerId:    getPeerID(),
// 		LocalIp:   getLocalIP(config.Config.Iface),
// 		Latitude:  latitude,
// 		Longitude: longitude,
// 		Protocol:  Protocol,
// 	}
// 	conn, err := grpc.Dial(RemoteProxy.Addr+":"+RemoteProxy.GrpcPort, grpc.WithInsecure())
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()
// 	client := proto.NewProxyClient(conn)

// 	// create tcp or udp connection. and set listen port number
// 	netConn, port, err := createNetConn()
// 	if err != nil {
// 		return err
// 	}
// 	NetConn = netConn
// 	request.LocalPort = port

// 	if _, err := client.PeerJoinRPC(context.Background(), request); err != nil {
// 		return err
// 	}

// 	// Set flag to decide finised the process of joining to Proxy
// 	IsCompletedJoinToProxy = true

// 	return nil
// }

// // return
// func createNetConn() (net.Conn, string, error) {
// 	switch Protocol {
// 	case "tcp":
// 		conn, err := net.Dial("tcp", RemoteProxy.Addr+":"+RemoteProxy.TcpPort)
// 		if err != nil {
// 			return nil, "", err
// 		}
// 		addr := conn.LocalAddr().String()
// 		port := extractPort(addr)
// 		return conn, port, nil
// 	case "udp":
// 		conn, err := net.Dial("udp", RemoteProxy.Addr+":"+RemoteProxy.UdpPort)
// 		if err != nil {
// 			return nil, "", err
// 		}
// 		addr := conn.LocalAddr().String()
// 		port := extractPort(addr)
// 		return conn, port, nil
// 	}
// 	return nil, "", nil
// }

// Only ipv4
func extractPort(addr string) string {
	return strings.Split(addr, ":")[1]
}
