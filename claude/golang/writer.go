package golang

import "net"

type Writer interface {
	Write(b []byte) (int, error)
}

type writer struct {
	conn net.Conn
}

func NewWriter(c net.Conn) Writer {
	return &writer{
		conn: c,
	}
}

func (w *writer) Write(b []byte) (int, error) {
	return w.conn.Write(b)
}
