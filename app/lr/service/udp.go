package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/common/proto"

	"github.com/Bo0km4n/claude/app/common/message"
	"github.com/Bo0km4n/claude/app/lr/config"
	"google.golang.org/grpc"
)

func ListenUDPBcastFromPeer() {
	addr := fmt.Sprintf("%s:%s", config.Config.UDP.Address, config.Config.UDP.Port)
	log.Printf("UDP Process is Running at %s\n", addr)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	inf, err := net.InterfaceByName("eth1")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenMulticastUDP("udp", inf, udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1500)
	for {
		length, remoteAddr, _ := conn.ReadFromUDP(buffer)
		request := &message.UDPBcastMessage{}
		if err := json.Unmarshal(buffer[:length], request); err != nil {
			log.Printf("Received from %v. Can't binding message to struct\n", remoteAddr)
			continue
		}

		log.Printf("Received from %v: %v\n", remoteAddr, request)
		request.ListenAddr = remoteAddr.IP.String()
		sendNoticeToPeer(request)
	}
}

func sendNoticeToPeer(m *message.UDPBcastMessage) {
	conn, err := grpc.Dial(m.ListenAddr+":"+m.ListenPort, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	client := proto.NewPeerClient(conn)
	if _, err := client.NoticeFromLRRPC(context.Background(), &proto.NoticeFromLRRequest{
		TcpPort:  config.Config.Claude.TcpPort,
		UdpPort:  config.Config.Claude.UdpPort,
		GrpcPort: config.Config.GRPC.Port,
		Addr:     config.Config.GRPC.Addr,
	}); err != nil {
		log.Println(err)
		return
	}
}
