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
	tcp      = flag.Bool("tcp", false, "use tcp")
	chunk    = flag.Int("chunk", 512, "chunk data size")
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

	log.Println("Type peer addr")
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

	writedSize := 0
	ticker := time.NewTicker(time.Minute * time.Duration(*minute))
	defer ticker.Stop()

	go func() {
		for {
			buf := make([]byte, *chunk)
			n, err := turnConn.Write(buf[:*chunk])
			if err != nil {
				log.Fatal(err)
			}
			writedSize += n
			fmt.Fprintf(os.Stdout, "\rwrite: %d", writedSize)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-ticker.C:
		log.Printf("Finished write data at: %d, size: %d\n", time.Now().UTC().UnixNano()/int64(time.Millisecond), writedSize)
	}
}
