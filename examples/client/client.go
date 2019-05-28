package main

import (
	"bufio"
	"fmt"
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
	repl(conn)
	<-quit
	// conn.SaveConnection()
	log.Println("exited")
}

func repl(conn *lib.Connection) {
	stdin := bufio.NewScanner(os.Stdin)
	fmt.Printf(">>> ")
	for stdin.Scan() {
		fmt.Printf(">>> ")
		text := stdin.Text()
		if text == "exit" {
			os.Exit(0)
		}
		if err := conn.Write([]byte(text)); err != nil {
			log.Fatal(err)
		}
	}
}
