package main

import (
	"github.com/Bo0km4n/claude/app/tablet/config"
	"github.com/Bo0km4n/claude/app/tablet/pkg/api"
	"github.com/Bo0km4n/claude/app/tablet/pkg/db"
)

func init() {
	config.InitConfig()
	db.InitMysql("claude")
	db.MigrateMysql()
}

func main() {
	api.GRPC()
}
