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
	Addr string
	Port string
}

var RemoteLR remoteLR

type PeerService struct{}

func (p *PeerService) NoticeFromLRRPC(ctx context.Context, in *proto.NoticeFromLRRequest) (*proto.Empty, error) {
	if in.Addr == "" || in.Port == "" {
		return nil, errors.New("LR information is invalid")
	}

	if RemoteLR.Addr != "" && RemoteLR.Port != "" {
		return &proto.Empty{}, nil
	}

	RemoteLR.Addr = in.Addr
	RemoteLR.Port = in.Port

	log.Printf("Registered LR | Addr: %s, Port: %s\n", RemoteLR.Addr, RemoteLR.Port)

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
		PeerId:    getPeerID(),
		LocalIp:   getLocalIP(config.Config.Iface),
		LocalPort: config.Config.Claude.Port,
		Latitude:  latitude,
		Longitude: longitude,
	}
	conn, err := grpc.Dial(RemoteLR.Addr+":"+RemoteLR.Port, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := proto.NewLRClient(conn)
	if _, err := client.PeerJoinRPC(context.Background(), request); err != nil {
		return err
	}
	return nil
}
