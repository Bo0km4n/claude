package service

import (
	"log"
	"net"
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/lr/config"
	"google.golang.org/grpc"
)

func LaunchService() {
	initService()
	initDaemon()
	go LRSvc.ListenUDPBcastFromPeer()
	go td.start()

	time.Sleep(2)
	if err := td.syncInit(); err != nil {
		log.Fatal(err)
	}
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
