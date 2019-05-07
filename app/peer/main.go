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
	go service.LaunchGRPCService(done, "udp")
	<-done

	time.Sleep(2)
	service.UDPBcast()
	for {
		if service.IsCompletedJoinToLR {
			return
		}
		time.Sleep(1)
	}
}
