package main

import (
	"github.com/Bo0km4n/claude/pkg/tablet/config"
	"github.com/Bo0km4n/claude/pkg/tablet/pkg/api"
	"github.com/Bo0km4n/claude/pkg/tablet/pkg/db"
)

func init() {
	config.InitConfig()
	db.InitMysql("claude")
	db.MigrateMysql()
}

func main() {
	api.GRPC()
}
