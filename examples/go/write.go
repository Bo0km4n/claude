package main

import (
	"log"
	"net"
	"os"

	"github.com/Bo0km4n/claude/claude/go/service"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Unexpected len(os.Args)=%d", len(os.Args))
	}
	seed := os.Args[1]
	service.SetProxyInformation(seed)

	proxyTcpAddr := service.GetProxyTCPAddr()
	conn, err := net.Dial("tcp", proxyTcpAddr)
	if err != nil {
		panic(err)
	}

	conn.Write([]byte(`hello world!`))
}
