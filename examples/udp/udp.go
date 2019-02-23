package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	addr := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
	conn, err := net.Dial("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("サーバへメッセージを送信.")
	conn.Write([]byte("Hello From Client."))

	fmt.Println("サーバからメッセージを受信。")
	buffer := make([]byte, 1500)
	length, _ := conn.Read(buffer)
	fmt.Printf("Receive: %s \n", string(buffer[:length]))
}
