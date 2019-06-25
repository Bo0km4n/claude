package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/Bo0km4n/claude/claude/golang/cio"
	"github.com/Bo0km4n/claude/claude/golang/service"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Arguments too short")
	}
	f, err := os.Open(os.Args[3])
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	seed := os.Args[1]
	service.SetProxyInformation(seed)

	proxyTcpAddr := service.GetProxyTCPAddr()
	conn, err := net.Dial("tcp", proxyTcpAddr)
	if err != nil {
		panic(err)
	}

	w := cio.NewWriter(conn)
	if _, err := w.Send(
		os.Args[2],
		data); err != nil {
		log.Fatal(err)
	}
}
