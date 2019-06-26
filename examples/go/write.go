package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/Bo0km4n/claude/claude/golang/cio"
	"github.com/Bo0km4n/claude/claude/golang/service"
)

var (
	useMulticast = flag.Bool("use_multicast", false, "use multicast udp")
	proxyAddr    = flag.String("proxy_addr", "", "proxy addr")
	iface        = flag.String("iface", "en0", "interface")
	seed         = flag.String("seed", "", "seed id")
)

func init() {
	flag.Parse()
}

func main() {
	service.SetProxyInformation(*seed, *iface, *proxyAddr, *useMulticast)

	proxyTcpAddr := service.GetProxyTCPAddr()
	conn, err := net.Dial("tcp", proxyTcpAddr)
	if err != nil {
		panic(err)
	}

	w := cio.NewWriter(conn)
	if _, err := w.Send(
		os.Args[2],
		[]byte(`hello world`)); err != nil {
		log.Fatal(err)
	}
}
