package service

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/peer/config"
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
