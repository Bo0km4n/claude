package service

import (
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/lr/config"
	"google.golang.org/grpc"
)

func LaunchService() {
	initService()
	go LRSvc.ListenUDPBcastFromPeer()
	launchPacketFilter()
	launchGRPCService()
}

func launchGRPCService() {
	port, err := net.Listen("tcp", ":"+config.Config.GRPC.Port)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()

	proto.RegisterLRServer(
		server,
		LRSvc,
	)

	log.Println("Start grpc services...")
	server.Serve(port)
}
