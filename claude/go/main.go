package main

import (
	"time"

	"github.com/Bo0km4n/claude/claude/go/config"
	"github.com/Bo0km4n/claude/claude/go/service"
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
		if service.IsCompletedJoinToProxy {
			return
		}
		time.Sleep(1)
	}
}
