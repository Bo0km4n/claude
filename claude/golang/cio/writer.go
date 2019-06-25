package cio

import (
	"net"

	"github.com/Bo0km4n/claude/claude/golang/packet"
	"github.com/Bo0km4n/claude/claude/golang/service"
	"github.com/k0kubun/pp"
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

func (w *writer) send(to string, body []byte) (int, error) {
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

func (w *writer) Send(to string, body []byte) (int, error) {
	if !(len(body) > packet.PACKET_SIZE-packet.HEADER_LENGTH) {
		return w.send(to, body)
	}

	offset := packet.PACKET_SIZE - packet.HEADER_LENGTH
	start := 0
	num := 0
	// split packet
	for {
		pp.Println(to, start, offset, offset-start)
		if offset > len(body) {
			offset = len(body)
			n, err := w.send(to, body[start:offset])
			if err != nil {
				return num, nil
			}
			num += n
			break
		}
		chunk := body[start:offset]
		n, err := w.send(to, chunk)
		if err != nil {
			return num, err
		}
		num += n
		start = offset
		offset += packet.PACKET_SIZE - packet.HEADER_LENGTH
	}
	return num, nil
}
