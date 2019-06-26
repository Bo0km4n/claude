package main

import (
	"log"
	"net"
	"os"

	"github.com/Bo0km4n/claude/claude/golang/cio"
	"github.com/Bo0km4n/claude/claude/golang/packet"
	"github.com/Bo0km4n/claude/claude/golang/service"
	"github.com/k0kubun/pp"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Unexpected len(os.Args)=%d", len(os.Args))
	}
	seed := os.Args[1]
	service.SetProxyInformation(seed, os.Args[2])

	proxyTcpAddr := service.GetProxyTCPAddr()
	conn, err := net.Dial("tcp", proxyTcpAddr)
	if err != nil {
		panic(err)
	}

	r := cio.NewReader(conn)
	for {
		buf := make([]byte, packet.PACKET_SIZE)
		n, err := r.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		pp.Println(n)
	}
}
