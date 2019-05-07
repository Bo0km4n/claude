package service

import (
	"crypto/sha256"
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/peer/config"
)

func getLocalIP(dev string) string {
	iface, err := net.InterfaceByName(dev)
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

	return addr.IP.String()
}

func GetPeerID() []byte {
	peerID := sha256.Sum256([]byte(config.Config.Claude.Credential))
	return peerID[:]
}
