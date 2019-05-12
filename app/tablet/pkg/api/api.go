package api

import (
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/tablet/config"
	"google.golang.org/grpc"
)

func GRPC() {
	port, err := net.Listen("tcp", ":"+config.Config.GRPC.Port)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()

	log.Println("Start grpc services...")
	server.Serve(port)
}
