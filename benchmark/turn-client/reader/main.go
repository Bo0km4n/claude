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
	tcp      = flag.Bool("tcp", false, "use tcp")
	chunk    = flag.Int("chunk", 512, "read buffer")
	minute   = flag.Int("minute", 5, "minute")
)

func init() {
	flag.Parse()
}

func main() {
	// Resolving to TURN server.
	var alloc *turnc.Allocation

	if *tcp {
		raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", *turnHost, *turnPort))
		if err != nil {
			panic(err)
		}
		c, err := net.DialTCP("tcp", nil, raddr)
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
		alloc = a
		log.Println("allocated relay addr:", a.Relayed().String())
	} else {
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
		alloc = a
		log.Println("allocated relay addr:", a.Relayed().String())
	}

	log.Println("Type peer address (IP:Port)")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	addr := scanner.Text()

	var turnConn *turnc.Connection

	if *tcp {
		peerAddr, resolveErr := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s", addr))
		if resolveErr != nil {
			panic(resolveErr)
		}
		permission, createErr := alloc.Create(peerAddr.IP)
		if createErr != nil {
			panic(createErr)
		}
		conn, err := permission.CreateTCP(peerAddr)
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		turnConn = conn
		log.Println("establish turn connection with TCP")
	} else {
		peerAddr, resolveErr := net.ResolveUDPAddr("udp", fmt.Sprintf("%s", addr))
		if resolveErr != nil {
			panic(resolveErr)
		}
		permission, createErr := alloc.Create(peerAddr.IP)
		if createErr != nil {
			panic(createErr)
		}
		conn, err := permission.CreateUDP(peerAddr)
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		turnConn = conn
		log.Println("establish turn connection with UDP")

	}
	// Connection implements net.Conn.
	// limit := *dataSize
	readSize := 0
	log.Println("Enter read loop")

	first := false
	ticker := time.NewTicker(time.Minute * time.Duration(*minute))
	go func() {
		for {
			buf := make([]byte, *chunk)
			n, err := turnConn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			if !first {
				first = true
				ticker = time.NewTicker(time.Minute * time.Duration(*minute))
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
