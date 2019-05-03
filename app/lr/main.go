package main

import (
	"github.com/Bo0km4n/claude/app/lr/config"
	"github.com/Bo0km4n/claude/app/lr/db"
	"github.com/Bo0km4n/claude/app/lr/service"
)

func init() {
	config.InitConfig()
	db.InitDB()
}

func main() {
	go service.ListenUDPBcastFromPeer()
	service.LaunchPacketFilter()
	service.LaunchGRPCService()
}
