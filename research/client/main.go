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
	// data, _ := ioutil.ReadFile("dummy.data")
	// conn.Write(data)
	for i := 0; i < 100; i++ {
		data := []byte(`claude_packet payload`)
		n, err := conn.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(n)
	}
}
