package service

import (
	"bytes"
	"crypto/sha256"
	"log"
	"net"
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
	var macAddr string
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				macAddr = i.HardwareAddr.String()
				break
			}
		}
	}
	peerID := sha256.Sum256([]byte(macAddr))
	return peerID[:]
}
