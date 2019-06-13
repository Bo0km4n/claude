package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:10000")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func(c net.Conn) {
			for {
				buf := make([]byte, 0xffff)
				n, err := conn.Read(buf)
				if err != nil && err != io.EOF {
					panic(err)
				}
				if err == io.EOF {
					break
				}
				log.Println(n, string(buf[:15]))
				conn.Write(buf[:n])
			}
		}(conn)
	}
}
