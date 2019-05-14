package main

import (
	"github.com/Bo0km4n/claude/app/lr/config"
	"github.com/Bo0km4n/claude/app/lr/repository"
	"github.com/Bo0km4n/claude/app/lr/service"
)

func init() {
	config.InitConfig()
	repository.InitDB()
}

func main() {
	service.LaunchService()
}
