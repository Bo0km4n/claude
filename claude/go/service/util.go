package service

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"log"
	"net"

	"github.com/Bo0km4n/claude/claude/go/config"
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

func getPeerID() []byte {
	proxyID := make([]byte, 4)
	id := sha256.Sum256([]byte(config.Config.Claude.Credential))
	binary.BigEndian.PutUint32(proxyID, RemoteProxy.ID)
	peerID := append([]byte{}, proxyID[:]...)
	peerID = append(peerID, id[:]...)
	return peerID[0:36]
}

func getPeerIDString() string {
	id := getPeerID()
	return base64.StdEncoding.EncodeToString(id)
}
