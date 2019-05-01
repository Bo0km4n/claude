package main

import (
	"github.com/Bo0km4n/claude/app/lr/config"
	"github.com/Bo0km4n/claude/app/lr/service"
)

func init() {
	config.InitConfig()
}

func main() {
	go service.ListenUDPBcastFromPeer()
	service.LaunchGRPCService()
}
