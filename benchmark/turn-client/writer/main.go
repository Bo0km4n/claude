package main

import (
	"bufio"
	"flag"
	"fmt"
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
	minute   = flag.Int("minute", 5, "minute")
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
	chunkSize := 1492
	writedSize := 0
	ticker := time.NewTicker(time.Minute * time.Duration(*minute))
	defer ticker.Stop()

	go func() {
		for {
			buf := make([]byte, chunkSize)
			n, err := conn.Write(buf[:chunkSize])
			if err != nil {
				log.Fatal(err)
			}
			writedSize += n
			fmt.Fprintf(os.Stdout, "\rwrite: %d", writedSize)
		}
	}()

	select {
	case <-ticker.C:
		log.Printf("Finished write data at: %d, size: %d\n", time.Now().UTC().UnixNano()/int64(time.Millisecond), writedSize)
	}
}
