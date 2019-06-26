package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/Bo0km4n/claude/claude/golang/cio"
	"github.com/Bo0km4n/claude/claude/golang/service"
)

var (
	useMulticast = flag.Bool("use_multicast", false, "use multicast udp")
	proxyAddr    = flag.String("proxy_addr", "", "proxy addr")
	iface        = flag.String("iface", "en0", "interface")
	seed         = flag.String("seed", "", "seed id")
	to           = flag.String("to", "", "remote peer id")
	file         = flag.String("file", "", "dummy data path")
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
	f, _ := os.Open(*file)
	data, _ := ioutil.ReadAll(f)
	log.Println("Start", time.Now().Unix())
	if _, err := w.Send(
		*to,
		data); err != nil {
		log.Fatal(err)
	}
}
