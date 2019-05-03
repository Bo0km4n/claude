package service

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/peer/config"
	"github.com/Bo0km4n/claude/app/peer/geo"
	"google.golang.org/grpc"
)

// global variable
type remoteLR struct {
	Addr     string
	GrpcPort string
	TcpPort  string
	UdpPort  string
}

var RemoteLR remoteLR
var IsCompletedJoinToLR bool

type PeerService struct{}

func (p *PeerService) NoticeFromLRRPC(ctx context.Context, in *proto.NoticeFromLRRequest) (*proto.Empty, error) {
	if in.Addr == "" || in.GrpcPort == "" || in.TcpPort == "" || in.UdpPort == "" {
		return nil, errors.New("LR information is invalid")
	}

	if RemoteLR.Addr != "" && RemoteLR.GrpcPort != "" && RemoteLR.TcpPort != "" && RemoteLR.UdpPort != "" {
		return &proto.Empty{}, nil
	}

	RemoteLR.Addr = in.Addr
	RemoteLR.GrpcPort = in.GrpcPort
	RemoteLR.TcpPort = in.TcpPort
	RemoteLR.UdpPort = in.UdpPort

	log.Printf("Registered LR | Addr: %s, GrpcPort: %s\n", RemoteLR.Addr, RemoteLR.GrpcPort)

	if err := peerJoin(); err != nil {
		log.Fatalf("Failed join to LR: %+v\n", err)
	}

	return &proto.Empty{}, nil
}

func LaunchGRPCService(done chan<- int) {
	port, err := net.Listen("tcp", ":"+config.Config.GRPC.Port)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()

	proto.RegisterPeerServer(
		server,
		&PeerService{},
	)

	log.Println("Start grpc services...")
	done <- 1
	server.Serve(port)
}

func peerJoin() error {
	latitude, longitude := geo.GetLocation()
	request := &proto.PeerJoinRequest{
		PeerId:    GetPeerID(),
		LocalIp:   getLocalIP(config.Config.Iface),
		LocalPort: config.Config.Claude.Port,
		Latitude:  latitude,
		Longitude: longitude,
	}
	conn, err := grpc.Dial(RemoteLR.Addr+":"+RemoteLR.GrpcPort, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := proto.NewLRClient(conn)
	if _, err := client.PeerJoinRPC(context.Background(), request); err != nil {
		return err
	}

	// Set flag to decide finised the process of joining to LR
	IsCompletedJoinToLR = true

	return nil
}
