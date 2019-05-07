package main

import (
	"crypto/sha256"
	"log"

	"github.com/Bo0km4n/claude/lib"
)

func main() {
	lib.InitEnv("udp")

	// dest PeerB = efgh
	// dest PeerA = abcd
	dest := sha256.Sum256([]byte(`efgh`))
	// dest := service.GetPeerID()
	conn, err := lib.NewConnection("udp", dest[:])
	if err != nil {
		log.Fatal(err)
	}
	// pp.Println(conn)
	conn.Ping()
}
