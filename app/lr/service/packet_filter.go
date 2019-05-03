package service

import (
	"fmt"
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/lr/config"
)

func LaunchPacketFilter() {

	// Debug tcp listener
	go func() {
		listen, err := net.Listen("tcp", ":"+config.Config.Claude.TcpPort)
		if err != nil {
			log.Fatal(err)
		}
		defer listen.Close()

		buf := make([]byte, 1024)
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
			}
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Received packets: %d\n", n)
		}
	}()

	// Debug udp listener
	go func() {
		conn, _ := net.ListenPacket("udp", ":"+config.Config.Claude.UdpPort)
		defer conn.Close()

		buffer := make([]byte, 1024)
		for {
			length, remoteAddr, _ := conn.ReadFrom(buffer)
			fmt.Printf("Received from %v: %v\n", remoteAddr, buffer[:length])
		}
	}()
}
