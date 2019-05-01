package main

import (
	"time"

	"github.com/Bo0km4n/claude/app/peer/config"
	"github.com/Bo0km4n/claude/app/peer/service"
)

func init() {
	config.InitConfig()
}

func main() {
	done := make(chan int)
	wait := make(chan struct{})

	go service.LaunchGRPCService(done)
	<-done

	time.Sleep(2)
	service.UDPBcast()
	<-wait
}
