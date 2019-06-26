package main

import (
	"log"
	"net"
	"time"
)

func main() {
	listen, _ := net.Listen("tcp", ":9000")
	for {
		conn, _ := listen.Accept()
		go func(c net.Conn) {
			limit := 1024000000
			readSize := 0
			for {
				buf := make([]byte, 0xffff)
				n, err := c.Read(buf)
				if err != nil {
					panic(err)
				}
				readSize += n
				if readSize >= limit {
					c.Close()
					break
				}
			}
			log.Println("Finished", time.Now().Unix())
			return
		}(conn)
	}
}