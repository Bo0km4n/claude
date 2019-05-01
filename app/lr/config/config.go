package config

import (
	"log"
	"net"

	"github.com/kelseyhightower/envconfig"
)

var Config *conf

type LR struct {
	Interface string `required:"true" default:"eth1"`
	GRPC      GRPC
	UDP       UDP
}

type GRPC struct {
	Addr string
	Port string `required:"true" default:"50051"`
}

type UDP struct {
	Address string `required:"true" default:"224.0.0.1"`
	Port    string `required:"true" default:"9000"`
}

type conf struct {
	Interface string
	GRPC      GRPC
	UDP       UDP
}

func InitConfig() {
	lr := LR{}
	if err := envconfig.Process("lr", &lr); err != nil {
		log.Fatal(err)
	}

	iface, err := net.InterfaceByName("eth1")
	if err != nil {
		log.Fatal(err)
	}
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		log.Fatal(err)
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}
					break
				}
			}
		}
	}

	lr.GRPC.Addr = addr.IP.String()

	Config = &conf{
		Interface: lr.Interface,
		GRPC:      lr.GRPC,
		UDP:       lr.UDP,
	}

}
