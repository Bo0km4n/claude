package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

var (
	to   = flag.String("to", "", "remote peer id")
	file = flag.String("file", "", "remote peer id")
)

func init() {
	flag.Parse()
}

func main() {
	conn, err := net.Dial("tcp", *to)
	if err != nil {
		panic(err)
	}
	f, _ := os.Open(*file)
	data, _ := ioutil.ReadAll(f)
	log.Println("Start", time.Now().Unix())
	_, err := conn.Write(data)
	if err != nil {
		panic(err)
	}
}
