package main

import (
	"log"
	"os"

	"github.com/Bo0km4n/claude/lib"
)

func main() {
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
	// pp.Println(conn)
	conn.Ping()
}
