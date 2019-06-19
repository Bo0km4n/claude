package repository

import (
	"github.com/Bo0km4n/claude/pkg/proxy/repository/pipe"
	"github.com/Bo0km4n/claude/pkg/proxy/repository/remotepeer"
)

func InitDB() {
	pipe.InitRepo()
	remotepeer.InitRepo()
}
