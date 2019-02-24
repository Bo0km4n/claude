package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/k0kubun/pp"
)

func main() {
	addr := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
	conn, err := net.Dial("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("サーバへメッセージを送信.")
	conn.Write(buildClaudeHeader())

	fmt.Println("サーバからメッセージを受信。")
	buffer := make([]byte, 1500)
	length, _ := conn.Read(buffer)
	fmt.Printf("Receive: %s \n", string(buffer[:length]))
}

func buildClaudeHeader() []byte {
	// SrcIP: 200.10.1.1
	// DstIP: 100.10.1.2
	// SrcPort: 50010
	// DstPort: 50011
	body := make([]byte, 16)
	msg := []byte("Hello world")
	srcIP := uint32(0xc80a0101)
	dstIP := uint32(0xa00a0102)
	srcPort := uint16(0xc35a)
	dstPort := uint16(0xc35b)
	checkSum := uint16(0x0101)
	length := uint16(16) + uint16(len(msg))

	binary.BigEndian.PutUint32(body, srcIP)
	binary.BigEndian.PutUint32(body[4:], dstIP)
	binary.BigEndian.PutUint16(body[8:], srcPort)
	binary.BigEndian.PutUint16(body[10:], dstPort)
	binary.BigEndian.PutUint16(body[12:], checkSum)
	binary.BigEndian.PutUint16(body[14:], length)

	body = append(body, msg...)
	pp.Println(body)
	return body
}
