package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	laddr := conn.LocalAddr()
	log.Println(laddr.String())
	// for i := 0; i < 6000; i++ {
	// 	var data []byte
	// 	if i == 0 {
	// 		data = []byte(`1 Init`)
	// 	} else {
	// 		data = []byte(`Hello world!`)
	// 	}
	// 	_, err := conn.Write(data)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// for {
	// 	buf := make([]byte, 0xffff)
	// 	n, err := conn.Read(buf)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println(n)
	// }
}
