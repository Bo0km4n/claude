package repository

import (
	"net"
	"sync"
)

// peerConnectionRepository has connections between Proxy and Peer
var peerConnectionRepository sync.Map

func RegisterPeerConnection(key string, protocol string, value net.Conn) {
	switch protocol {
	case "tcp":
		conn := value.(*net.TCPConn)
		peerConnectionRepository.Store(key, conn)
	case "udp":
		conn := value.(*net.UDPConn)
		peerConnectionRepository.Store(key, conn)
	}
}

func LoadPeerConnection(key, protocol string) (net.Conn, bool) {
	v, ok := peerConnectionRepository.Load(key)
	if !ok {
		return nil, false
	}
	return v.(net.Conn), ok
}

// proxyConnectionRepository has connections between itself and remote Proxy
var proxyConnectionRepository sync.Map
