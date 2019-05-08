package main

import (
	"crypto/sha256"
	"log"
	"os"

	"github.com/Bo0km4n/claude/lib"
)

func main() {
	lib.InitEnv(os.Args[1])

	// dest PeerB = efgh
	// dest PeerA = abcd

	destID := os.Args[2]
	dest := sha256.Sum256([]byte(destID))
	// dest := service.GetPeerID()
	conn, err := lib.NewConnection(os.Args[1], dest[:])
	if err != nil {
		log.Fatal(err)
	}
	// pp.Println(conn)
	conn.Ping()
}
