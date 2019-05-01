package service

import (
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/lr/config"
	"google.golang.org/grpc"
)

func LaunchGRPCService() {
	port, err := net.Listen("tcp", ":"+config.Config.GRPC.Port)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()

	proto.RegisterLRServer(
		server,
		&LRService{},
	)

	log.Println("Start grpc services...")
	server.Serve(port)
}
