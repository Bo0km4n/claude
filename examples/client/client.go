package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Bo0km4n/claude/lib"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	lib.InitConfig()
	lib.ConnectToLR(os.Args[1])

	// dest PeerB = efgh
	// dest PeerA = abcd

	dest := lib.DeserializeID(os.Args[2])
	// dest := service.GetPeerID()
	conn, err := lib.NewConnection(os.Args[1], dest[:])
	if err != nil {
		log.Fatal(err)
	}
	go conn.Ping()
	<-quit
	conn.SaveConnection()
	log.Println("exited")

	// remoteAddr, _ := net.ResolveUDPAddr("udp", "192.168.10.100:9611")

	// conn, err := net.DialUDP("udp", &net.UDPAddr{
	// 	IP: net.IP{
	// 		0xc0, 0xa8, 0x0a, 0x65,
	// 	},
	// 	Port: 57254,
	// 	Zone: "",
	// }, remoteAddr)
	// if err != nil {
	// 	panic(err)
	// }

	// conn.Write([]byte("Hello world"))

}
