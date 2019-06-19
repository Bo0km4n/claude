package cio

import "net"

type ReaderAndWriter interface {
	Reader
	Writer
}

type readerAndWriter struct {
	conn net.Conn
}
