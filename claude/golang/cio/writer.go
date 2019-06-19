package cio

import (
	"net"

	"github.com/Bo0km4n/claude/claude/golang/packet"
	"github.com/Bo0km4n/claude/claude/golang/service"
)

type Writer interface {
	Send(to string, b []byte) (int, error)
}

type writer struct {
	conn net.Conn
	ps   *service.PeerService
}

func NewWriter(c net.Conn) Writer {
	return &writer{
		ps:   service.PeerSvc,
		conn: c,
	}
}

func (w *writer) write(b []byte) (int, error) {
	return w.conn.Write(b)
}

func (w *writer) Send(to string, body []byte) (int, error) {
	newPacket := packet.GeneratePacket()
	if err := newPacket.SetToID(to); err != nil {
		return 0, err
	}
	if err := newPacket.SetFromID(w.ps.ID); err != nil {
		return 0, err
	}
	if err := newPacket.SetPayload(body); err != nil {
		return 0, err
	}
	msg := newPacket.Serialize()
	return w.conn.Write(msg)
}
