package lib

import (
	"errors"
	"net"

	"github.com/Bo0km4n/claude/app/peer/service"
)

type Connection struct {
	NetConn           net.Conn
	Protocol          string
	DestinationPeerID []byte
	SourcePeerID      []byte
}

func NewConnection(protocol string, dest []byte) (*Connection, error) {
	switch protocol {
	case "udp":
		conn, err := net.Dial("udp", service.RemoteLR.Addr+":"+service.RemoteLR.UdpPort)
		if err != nil {
			return nil, err
		}

		return &Connection{
			NetConn:           conn,
			Protocol:          "udp",
			DestinationPeerID: dest,
			SourcePeerID:      service.GetPeerID(),
		}, nil
	case "tcp":
		conn, err := net.Dial("tcp", service.RemoteLR.Addr+":"+service.RemoteLR.TcpPort)
		if err != nil {
			return nil, err
		}

		return &Connection{
			NetConn:           conn,
			Protocol:          "tcp",
			DestinationPeerID: dest,
			SourcePeerID:      service.GetPeerID(),
		}, nil
	}
	return nil, errors.New("Not found network")
}
