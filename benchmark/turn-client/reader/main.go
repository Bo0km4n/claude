package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"flag"

	"time"

	"gortc.io/turnc"
)

var (
	turnHost = flag.String("turn", "127.0.0.1", "turn server addr")
	turnPort = flag.String("p", "9610", "turn server port")
	dd       = flag.String("dd", "./dd.data", "dummy data file path")
	dataSize = flag.Int("ds", 1, "1 * GigaByte")
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
	// peerAddr, resolveErr := net.ResolveUDPAddr("udp", "10.0.0.1:34587")
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
	// limit := *dataSize
	readSize := 0
	log.Println("Enter read loop")

	first := false
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for {
			buf := make([]byte, 1492)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			if !first {
				first = true
				ticker = time.NewTicker(time.Minute * 5)
				defer ticker.Stop()
			}
			readSize += n
			fmt.Fprintf(os.Stdout, "\rread size: %d", readSize)
		}
	}()

	select {
	case <-ticker.C:
		log.Println("Finished", time.Now().UTC().UnixNano()/int64(time.Millisecond), readSize)
	}
}
