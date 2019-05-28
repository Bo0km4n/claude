package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Bo0km4n/claude/lib"
	"github.com/k0kubun/pp"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	lib.InitConfig()
	lib.ConnectToLR(os.Args[1])
	entries, err := lib.LookUpPeers(35.681236, 139.767125, 10.0)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(entries)
}
