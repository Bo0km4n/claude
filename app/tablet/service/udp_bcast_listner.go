package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type UDPBcastMessage struct {
	ListenPort string `json:"listen_port"`
}

func ListenUDPBcastFromLR() {
	fmt.Println("Server is Running at 224.0.0.1:9000")
	udpAddr, err := net.ResolveUDPAddr("udp", "224.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	inf, err := net.InterfaceByName("eth1")
	if err != nil {
		panic(err)
	}

	// set sock option

	conn, err := net.ListenMulticastUDP("udp", inf, udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1500)
	for {
		// 通信読込 + 接続相手アドレス情報が受取
		length, remoteAddr, _ := conn.ReadFromUDP(buffer)
		request := &UDPBcastMessage{}
		if err := json.Unmarshal(buffer[:length], request); err != nil {
			log.Printf("Received from %v. Can't binding message to struct\n", remoteAddr)
			continue
		}

		fmt.Printf("Received from %v: %v\n", remoteAddr, request)
	}
}
