package main

import (
	"log"

	peerService "github.com/Bo0km4n/claude/app/peer/service"
	"github.com/Bo0km4n/claude/lib"
)

func main() {
	lib.InitEnv()
	dest := peerService.GetPeerID()
	conn, err := lib.NewConnection("udp", dest)
	if err != nil {
		log.Fatal(err)
	}
	// pp.Println(conn)
	conn.Ping()
}
