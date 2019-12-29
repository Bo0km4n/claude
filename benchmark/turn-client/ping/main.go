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
	turnHost  = flag.String("turn", "127.0.0.1", "turn server addr")
	turnPort  = flag.String("p", "9610", "turn server port")
	tcp       = flag.Bool("tcp", false, "use tcp")
	frequency = flag.Int("f", 1, "second")
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
		permission, createErr := a.Create(peerAddr.IP)
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
		permission, createErr := a.Create(peerAddr.IP)
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
	ticker := time.NewTicker(time.Second * time.Duration(*frequency))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := turnConn.Write([]byte("ping"))
			if err != nil {
				log.Fatal(err)
			}
			log.Println("ping")
		}
	}
}
