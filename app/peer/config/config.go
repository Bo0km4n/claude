package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var Config *conf

type Peer struct {
	GRPC GRPC
}

type GRPC struct {
	Port string `required:"true" default:"50051"`
}

type conf struct {
	GRPC GRPC
}

func InitConfig() {
	peer := Peer{}
	if err := envconfig.Process("peer", &peer); err != nil {
		log.Fatal(err)
	}

	Config = &conf{
		GRPC: peer.GRPC,
	}
}
