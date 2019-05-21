package service

import (
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/lr/config"
	"google.golang.org/grpc"
)

func LaunchService() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	initService()
	initDaemon()
	go LRSvc.ListenUDPBcastFromPeer()
	go td.start()

	launchPacketFilter()
	go launchGRPCService()

	time.Sleep(2)
	if err := td.syncInit(); err != nil {
		log.Fatal(err)
	}

	<-quit
	log.Println("Interrupted LR Server")
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
