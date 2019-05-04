package main

import (
	"crypto/sha256"
	"log"

	"github.com/Bo0km4n/claude/lib"
)

func main() {
	lib.InitEnv()
	dest := sha256.Sum256([]byte(`hoge`))
	conn, err := lib.NewConnection("tcp", dest[:])
	if err != nil {
		log.Fatal(err)
	}
	// pp.Println(conn)
	conn.Ping()
}
