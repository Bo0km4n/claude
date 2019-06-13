package main

import (
	"log"
	"net"

	"golang.org/x/sync/errgroup"
)

func closeConn(in <-chan *net.TCPConn) {
	for conn := range in {
		conn.Close()
	}
}

func handleConn(in *net.TCPConn) {
	defer in.Close()
	log.Println(in.RemoteAddr().String())
	rConn, err := net.Dial("tcp", "localhost:10000")
	if err != nil {
		panic(err)
	}
	log.Println("Handle new connection")
	rTcpConn, _ := rConn.(*net.TCPConn)
	proxyConn(in, rTcpConn)
}

func proxyConn(client, server *net.TCPConn) {
	defer client.Close()
	defer server.Close()

	var eg errgroup.Group

	eg.Go(func() error { return relay(client, server) })
	eg.Go(func() error { return relay(server, client) })
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	log.Println("Finished relay")
}

func relay(from, to *net.TCPConn) error {
	buff := make([]byte, 0xffff)
	for {
		n, err := from.Read(buff)
		if err != nil {
			return err
		}
		b := buff[:n]

		header := []byte(`HEADER`)
		body := append(header, b[:]...)

		n, err = to.Write(body)
		if err != nil {
			return err
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		tcpConn, _ := conn.(*net.TCPConn)
		go handleConn(tcpConn)
	}
}
