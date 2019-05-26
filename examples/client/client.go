package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Bo0km4n/claude/lib"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	lib.InitConfig()
	lib.ConnectToLR(os.Args[1])

	dest := lib.DeserializeID(os.Args[2])
	conn, err := lib.NewConnection(os.Args[1], dest[:])
	if err != nil {
		log.Fatal(err)
	}
	conn.RegisterHandler(
		func(b []byte) error {
			log.Println(string(b))
			return nil
		})
	go conn.Serve()
	<-quit
	// conn.SaveConnection()
	log.Println("exited")
}
