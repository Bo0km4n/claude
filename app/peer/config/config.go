package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var Config *conf

type Peer struct {
	Iface       string `required:"true" default:"eth1"`
	GRPC        GRPC
	Claude      Claude
	NetConnFile string `required:"true" envconfig:"NET_CONN_FILE"`
}

type GRPC struct {
	Port string `required:"true" default:"50051"`
}

type Claude struct {
	Credential string `required:"true"`
	Port       string `required:"true" default:"9610"`
}

type conf struct {
	Iface       string
	GRPC        GRPC
	Claude      Claude
	NetConnFile string
}

func InitConfig() {
	peer := Peer{}
	if err := envconfig.Process("peer", &peer); err != nil {
		log.Fatal(err)
	}

	Config = &conf{
		Iface:       peer.Iface,
		GRPC:        peer.GRPC,
		Claude:      peer.Claude,
		NetConnFile: peer.NetConnFile,
	}
}
