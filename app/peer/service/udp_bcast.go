package service

import (
	"encoding/json"
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/common/message"
	"github.com/Bo0km4n/claude/app/peer/config"
)

func UDPBcast() {
	remoteUDPAddr, err := net.ResolveUDPAddr("udp", "224.0.0.1:9000")
	if err != nil {
		log.Fatal(err)
	}

	// Get local eth1 address
	ief, err := net.InterfaceByName("eth1")
	if err != nil {
		log.Fatal(err)
	}
	addrs, err := ief.Addrs()
	if err != nil {
		log.Fatal(err)
	}

	localUDPAddr := &net.UDPAddr{
		IP: addrs[0].(*net.IPNet).IP,
	}

	conn, err := net.DialUDP("udp", localUDPAddr, remoteUDPAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	msg, err := json.Marshal(buildRequest())
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Send multicast")
}

func buildRequest() *message.UDPBcastMessage {
	return &message.UDPBcastMessage{
		ListenPort: config.Config.GRPC.Port,
	}
}
