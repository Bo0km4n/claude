package lib

import (
	"time"

	"github.com/Bo0km4n/claude/app/peer/service"
)

func SetLR(protocol string) {
	done := make(chan int)

	go service.LaunchGRPCService(done, protocol)
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
