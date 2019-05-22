package service

import (
	"context"
	"errors"
	"log"
	"net"
	"strings"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/peer/config"
	"github.com/Bo0km4n/claude/app/peer/geo"
	"google.golang.org/grpc"
)

// global variable
type remoteLR struct {
	ID       uint32
	Addr     string
	GrpcPort string
	TcpPort  string
	UdpPort  string
}

var Protocol string
var NetConn net.Conn
var RemoteLR remoteLR
var IsCompletedJoinToLR bool
var PeerSvc *PeerService

type PeerService struct {
	Seed string
	ID   string
}

func (p *PeerService) NoticeFromLRRPC(ctx context.Context, in *proto.NoticeFromLRRequest) (*proto.Empty, error) {
	if in.Addr == "" || in.GrpcPort == "" || in.TcpPort == "" || in.UdpPort == "" {
		return nil, errors.New("LR information is invalid")
	}

	if RemoteLR.Addr != "" && RemoteLR.GrpcPort != "" && RemoteLR.TcpPort != "" && RemoteLR.UdpPort != "" {
		return &proto.Empty{}, nil
	}

	RemoteLR.ID = in.Id
	RemoteLR.Addr = in.Addr
	RemoteLR.GrpcPort = in.GrpcPort
	RemoteLR.TcpPort = in.TcpPort
	RemoteLR.UdpPort = in.UdpPort

	log.Printf("Registered LR | Addr: %s, GrpcPort: %s\n", RemoteLR.Addr, RemoteLR.GrpcPort)

	if err := peerJoin(); err != nil {
		log.Fatalf("Failed join to LR: %+v\n", err)
	}

	// set peer id
	p.Seed = config.Config.Claude.Credential
	p.ID = getPeerIDString()

	return &proto.Empty{}, nil
}

func LaunchGRPCService(done chan<- int, protocol string) {
	Protocol = protocol
	port, err := net.Listen("tcp", ":"+config.Config.GRPC.Port)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	PeerSvc = &PeerService{}
	proto.RegisterPeerServer(
		server,
		PeerSvc,
	)

	log.Println("Start grpc services...")
	done <- 1
	server.Serve(port)
}

func peerJoin() error {
	latitude, longitude := geo.GetLocation()
	request := &proto.PeerJoinRequest{
		PeerId:    getPeerID(),
		LocalIp:   getLocalIP(config.Config.Iface),
		Latitude:  latitude,
		Longitude: longitude,
		Protocol:  Protocol,
	}
	conn, err := grpc.Dial(RemoteLR.Addr+":"+RemoteLR.GrpcPort, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := proto.NewLRClient(conn)

	// create tcp or udp connection. and set listen port number
	netConn, port, err := createNetConn()
	if err != nil {
		return err
	}
	NetConn = netConn
	request.LocalPort = port

	if _, err := client.PeerJoinRPC(context.Background(), request); err != nil {
		return err
	}

	// Set flag to decide finised the process of joining to LR
	IsCompletedJoinToLR = true

	return nil
}

// return
func createNetConn() (net.Conn, string, error) {
	switch Protocol {
	case "tcp":
		conn, err := net.Dial("tcp", RemoteLR.Addr+":"+RemoteLR.TcpPort)
		if err != nil {
			return nil, "", err
		}
		addr := conn.LocalAddr().String()
		port := extractPort(addr)
		return conn, port, nil
	case "udp":
		conn, err := net.Dial("udp", RemoteLR.Addr+":"+RemoteLR.UdpPort)
		if err != nil {
			return nil, "", err
		}
		addr := conn.LocalAddr().String()
		port := extractPort(addr)
		return conn, port, nil
	}
	return nil, "", nil
}

// Only ipv4
func extractPort(addr string) string {
	return strings.Split(addr, ":")[1]
}
