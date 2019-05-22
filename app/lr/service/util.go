package service

import (
	"bytes"
	"net"
)

func getMacAddr(iframe string) (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Name == iframe && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}
