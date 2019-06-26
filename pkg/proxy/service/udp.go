package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/Bo0km4n/claude/pkg/common/proto"

	"github.com/Bo0km4n/claude/pkg/common/message"
	"github.com/Bo0km4n/claude/pkg/proxy/config"
	"google.golang.org/grpc"
)

func (proxys *ProxyService) ListenUDPFromPeer(isMulticast bool) {
	addr := fmt.Sprintf("%s:%s", config.Config.UDP.Address, config.Config.UDP.Port)
	log.Printf("UDP Process is Running at %s\n", addr)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	inf, err := net.InterfaceByName(config.Config.Interface)
	if err != nil {
		panic(err)
	}

	var conn *net.UDPConn
	if isMulticast {
		udpConn, err := net.ListenMulticastUDP("udp", inf, udpAddr)
		if err != nil {
			panic(err)
		}
		conn = udpConn
	} else {
		udpConn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			panic(err)
		}
		conn = udpConn
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
		proxys.sendNoticeToPeer(request)
	}
}

func (proxys *ProxyService) sendNoticeToPeer(m *message.UDPBcastMessage) {
	conn, err := grpc.Dial(m.ListenAddr+":"+m.ListenPort, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	client := proto.NewPeerClient(conn)
	if _, err := client.NoticeFromProxyRPC(context.Background(), &proto.NoticeFromProxyRequest{
		Id:       proxys.ID,
		TcpPort:  config.Config.Claude.UpTcpPort,
		UdpPort:  config.Config.Claude.UpUdpPort,
		GrpcPort: config.Config.GRPC.Port,
		Addr:     config.Config.GRPC.Addr,
	}); err != nil {
		log.Println(err)
		return
	}
}
