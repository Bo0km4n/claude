package golang

import "net"

type ReaderAndWriter interface {
	Read(b []byte) (int, error)
	Write(b []byte) (int, error)
}

type readerAndWriter struct {
	conn net.Conn
}

func (rw *readerAndWriter) Read(b []byte) (int, error) {
	return rw.conn.Read(b)
}

func (rw *readerAndWriter) Write(b []byte) (int, error) {
	return rw.conn.Write(b)
}
