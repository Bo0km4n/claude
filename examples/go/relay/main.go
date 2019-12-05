package main

import (
	"flag"
	"fmt"
	"net"

	"log"
	"strings"

	"bufio"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

var (
	mod  = flag.String("m", "server", "mode")
	host = flag.String("h", "127.0.0.1:8080", "remote host")
	name = flag.String("n", "a", "server name")
)

func init() {
	flag.Parse()
}

func main() {
	switch *mod {
	case "server":
		fmt.Println("server")
		server()
	case "client":
		fmt.Println("client")
		client()
	}
}

func server() {
	listen, err := net.Listen("tcp", "0.0.0.0:9610")
	if err != nil {
		log.Fatal("tcp://0.0.0.0:9610のリッスンに失敗しました")
	}
	fmt.Println("0.0.0.0:9610で受付開始しました")

	// コネクションを受け付ける
	for {
		conn, err := listen.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatal("コネクションを確立できませんでした")
		}
		// リクエスト元のアドレス
		remoteAddr := conn.RemoteAddr()
		fmt.Printf("[Remote Address]\n%s\n", remoteAddr)
		addrs := strings.Split(remoteAddr.String(), ":")

		conn.Write([]byte(fmt.Sprintf("%s:%s\n", addrs[0], addrs[1])))
		pp.Printf("writed to %s", remoteAddr.String())
	}

	// コネクションを切断する
}

type BindResponse struct {
	RemoteIP   string `json:"remote_ip"`
	RemotePort string `json:"remote_port"`
}

func client() {
	conn, err := net.Dial("tcp", *host)
	if err != nil {
		log.Fatal(err)
	}
	br := bufio.NewReader(conn)
	l, _, err := br.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	translatedPort := strings.Split(string(l), ":")[1]
	pp.Println(string(l), conn.LocalAddr().String())
	localAddrs := strings.Split(conn.LocalAddr().String(), ":")
	bindedLocalPort := localAddrs[1]

	conn.Close()
	listenRelayServer(bindedLocalPort, translatedPort)
}

func listenRelayServer(port, translatedPort string) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		log.Printf("[PING:Remote] %s", c.Request.RemoteAddr)
		c.JSON(200, gin.H{"message": "ok"})
	})

	log.Printf("Local Port: %s <---> Translated Port: %s", port, translatedPort)
	r.Run(fmt.Sprintf(":%s", port))
}
