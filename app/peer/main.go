package main

import (
	"log"
	"net"
)

func main() {
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

	_, err = conn.Write([]byte(`{"port":"6000"}`))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Send multicast")
}
