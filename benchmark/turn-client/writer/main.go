package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"time"

	"gortc.io/turnc"
)

var (
	turnHost = flag.String("turn", "127.0.0.1", "turn server addr")
	turnPort = flag.String("p", "9610", "turn server port")
	dd       = flag.String("dd", "./dd.data", "dummy data file path")
)

func init() {
	flag.Parse()
}

func main() {
	// Resolving to TURN server.
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", *turnHost, *turnPort))
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

	log.Println("Type peer port")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	portStr := scanner.Text()

	peerAddr, resolveErr := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%s", portStr))
	if resolveErr != nil {
		panic(resolveErr)
	}
	permission, createErr := a.Create(peerAddr.IP)
	if createErr != nil {
		panic(createErr)
	}
	conn, err := permission.CreateUDP(peerAddr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	// Connection implements net.Conn.
	f, err := os.Open(*dd)
	if err != nil {
		log.Fatal(err)
	}

	chunkSize := 1492
	for {
		buf := make([]byte, chunkSize)
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Fatal(err)
		}
		log.Println("write", n)
	}
	log.Printf("Finished write data at: %d\n", time.Now().UTC().UnixNano()/int64(time.Millisecond))

	time.Sleep(1000)
}