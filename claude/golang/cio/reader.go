package cio

import (
	"net"
)

type Reader interface {
	Read(b []byte) (int, error)
}

type reader struct {
	conn net.Conn
}

func NewReader(c net.Conn) Reader {
	return &reader{
		conn: c,
	}
}

func (w *reader) Read(b []byte) (int, error) {
	return w.conn.Read(b)
}
