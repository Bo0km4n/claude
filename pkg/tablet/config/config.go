package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var Config *conf

type Tablet struct {
	GRPC GRPC
}

type GRPC struct {
	Port string `required:"true" default:"50051"`
}

type conf struct {
	GRPC GRPC
}

func InitConfig() {
	tablet := Tablet{}
	if err := envconfig.Process("tablet", &tablet); err != nil {
		log.Fatal(err)
	}

	Config = &conf{
		GRPC: tablet.GRPC,
	}
}
