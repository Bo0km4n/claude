package main

import (
	"log"
	"os"
	"time"

	"github.com/Bo0km4n/claude/lib"
)

var data []byte
var counter int

func main() {
	lib.InitConfig()
	lib.ConnectToLR(os.Args[1])
	quit := make(chan struct{})
	dest := lib.DeserializeID(os.Args[2])

	// Wait another peer connect
	time.Sleep(5)

	conn, err := lib.NewConnection(os.Args[1], dest[:])
	if err != nil {
		log.Fatal(err)
	}
	conn.RegisterHandler(
		func(c *lib.Connection, b []byte) error {
			c.Write(b)
			counter++
			log.Println(counter)
			if counter >= 100 {
				quit <- struct{}{}
			}
			return nil
		})
	go conn.Serve()

	if os.Getenv("START") == "on" {
		start(conn)
	}

	<-quit
}

func start(conn *lib.Connection) {
	conn.Write(data)
}

func calcElpasedTime(before time.Time) (time.Time, time.Duration) {
	now := time.Now()
	duration := now.Sub(before)
	return now, duration
}
