package service

import (
	"log"
	"net"
	"time"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"github.com/Bo0km4n/claude/pkg/proxy/config"
	"github.com/Bo0km4n/claude/pkg/proxy/proxy/tcp"
	"google.golang.org/grpc"
)

func LaunchService() {
	initService()
	initDaemon()
	go ProxySvc.ListenUDPBcastFromPeer()
	go td.start()

	go launchGRPCService()

	time.Sleep(2)
	if err := td.syncInit(); err != nil {
		log.Fatal(err)
	}
	tcp.NewProxy().Serve()
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
