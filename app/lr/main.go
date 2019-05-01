package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
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
	remoteUDPAddr, err := net.ResolveUDPAddr("udp", "224.0.0.1:9000")
	conn, err := net.DialUDP("udp", localUDPAddr, remoteUDPAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("サーバへメッセージを送信.")

	jsonMsg := `{
        "listen_port": "6060"
    }`
	conn.Write([]byte(jsonMsg))
	fmt.Println("done")
}
