package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"flag"

	"gortc.io/turnc"
)

var (
	turnServer = flag.String("turn", "127.0.0.1:3678", "turn server addr")
	dd         = flag.String("dd", "./dd.data", "dummy data file path")
)

func init() {
	flag.Parse()
}

func main() {
	// Resolving to TURN server.
	raddr, err := net.ResolveUDPAddr("udp", "10.128.0.8:9610")
	if err != nil {
		panic(err)
	}
	c, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		panic(err)
	}
	client, clientErr := turnc.New(turnc.Options{
		Conn: c,
		// Credentials:
		Username: "user1",
		Password: "pass1",
	})
	if clientErr != nil {
		panic(clientErr)
	}
	a, allocErr := client.Allocate()
	if allocErr != nil {
		panic(allocErr)
	}
	log.Println("allocated relay addr:", a.Relayed().String())

	log.Println("Type peer ip:port")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	addrStr := scanner.Text()

	peerAddr, resolveErr := net.ResolveUDPAddr("udp", addrStr)
	if resolveErr != nil {
		panic(resolveErr)
	}
	permission, createErr := a.Create(peerAddr.IP)
	if createErr != nil {
		panic(createErr)
	}
	conn, err := permission.CreateUDP(peerAddr)
	if err != nil {
		panic(err)
	}
	// Connection implements net.Conn.
	if _, writeRrr := fmt.Fprint(conn, "hello world!"); writeRrr != nil {
		panic(writeRrr)
	}
	buf := make([]byte, 1500)
	n, readErr := conn.Read(buf)
	if readErr != nil {
		panic(readErr)
	}
	log.Printf("got message: %s", string(buf[:n]))
	// Also you can use ChannelData messages to reduce overhead:
	if err := conn.Bind(); err != nil {
		panic(err)
	}
}
