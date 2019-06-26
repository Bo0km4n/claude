package main

import (
	"flag"
	"log"
	"net"

	"github.com/Bo0km4n/claude/claude/golang/cio"
	"github.com/Bo0km4n/claude/claude/golang/packet"
	"github.com/Bo0km4n/claude/claude/golang/service"
	"github.com/k0kubun/pp"
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
