package config

import (
	"log"
	"net"

	"github.com/kelseyhightower/envconfig"
)

var Config *conf

type PROXY struct {
	Interface string `required:"true" default:"eth1"`
	Claude    Claude
	GRPC      GRPC
	UDP       UDP
	Tablet    Tablet
}

type GRPC struct {
	Addr string
	Port string `required:"true" default:"50051"`
}

type UDP struct {
	Address string `required:"true" default:"224.0.0.1"`
	Port    string `required:"true" default:"9000"`
}

type Claude struct {
	UniqueKey   string
	UpTcpPort   string `required:"true" default:"9610"`
	DownTcpPort string `required:"true" default:"9611"`
	UpUdpPort   string `required:"true" default:"8610"`
	DownUdpPort string `required:"true" default:"8611"`
}

type Tablet struct {
	IP   string
	Port string
}

type conf struct {
	Claude    Claude
	Interface string
	GRPC      GRPC
	UDP       UDP
	Tablet    Tablet
}

func InitConfig() {
	proxy := PROXY{}
	if err := envconfig.Process("proxy", &proxy); err != nil {
		log.Fatal(err)
	}

	iface, err := net.InterfaceByName(proxy.Interface)
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

	proxy.GRPC.Addr = addr.IP.String()

	Config = &conf{
		Claude:    proxy.Claude,
		Interface: proxy.Interface,
		GRPC:      proxy.GRPC,
		UDP:       proxy.UDP,
		Tablet:    proxy.Tablet,
	}

}
