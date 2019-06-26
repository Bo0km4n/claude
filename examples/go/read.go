package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/Bo0km4n/claude/claude/golang/cio"
	"github.com/Bo0km4n/claude/claude/golang/packet"
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

	r := cio.NewReader(conn)
	limit := 1024000000
	readSize := 0
	for {
		buf := make([]byte, packet.PACKET_SIZE)
		n, err := r.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		readSize += n
		if readSize >= limit {
			conn.Close()
			break
		}
	}

	log.Println("Finished", time.Now().UTC().UnixNano())
}
