package service

import (
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"github.com/Bo0km4n/claude/pkg/proxy/config"
	"google.golang.org/grpc"
)

func LaunchService() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	initService()
	initDaemon()
	go ProxySvc.ListenUDPBcastFromPeer()
	go td.start()

	go launchGRPCService()

	time.Sleep(2)
	if err := td.syncInit(); err != nil {
		log.Fatal(err)
	}

	<-quit
	log.Println("Interrupted Proxy Server")
}

func launchGRPCService() {
	port, err := net.Listen("tcp", ":"+config.Config.GRPC.Port)
	if err != nil {
		log.Fatal("launchGRPCService: ", err)
	}
	server := grpc.NewServer()

	proto.RegisterProxyServer(
		server,
		ProxySvc,
	)

	log.Println("Start grpc services...")
	server.Serve(port)
}
