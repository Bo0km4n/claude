package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/Bo0km4n/claude/lib"
	"github.com/k0kubun/pp"
)

var data []byte
var counter int
var before time.Time

func main() {
	withProxy()
	// withoutProxy()
}

func withoutProxy() {
	mode := os.Getenv("MODE")
	switch mode {
	case "SERVER":
		withoutProxyServer()
	case "CLIENT":
		withoutProxyClient()
	}
}

func withoutProxyServer() {
	listener, error := net.Listen("tcp", "192.168.10.101:10000")
	buf := make([]byte, 1024)
	if error != nil {
		panic(error)
	}
	for {
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		conn.Read(buf)
		pp.Println("READ")
		conn.Write(buf)
	}
}

func withoutProxyClient() {
	conn, err := net.Dial("tcp", "192.168.10.101:10000")
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)
	for {
		conn.Write(data)
		conn.Read(buf)
	}
}

var (
	Mode     = flag.String("m", "server", "application mode")
	DataPath = flag.String("d", "data path", "dummy data file")
)

func withProxy() {
	lib.InitConfig()
	lib.ConnectToProxy(os.Args[1])

	switch *Mode {
	case "server":
		quit := make(chan struct{})
		// Wait another peer connect
		time.Sleep(5)

		conn, err := lib.NewConnection(os.Args[1], dest[:])
		if err != nil {
			log.Fatal(err)
		}
		conn.RegisterHandler(
			func(c *lib.Connection, b []byte) error {
				c.Write(b)
				counter++
				if counter == 1 {
					initTime()
				} else {
					newTime, elapsed := calcElpasedTime(before)
					before = newTime
					log.Printf("[Term %d] Elapsed Time: %v\n", counter, elapsed)
				}
				if counter >= 100 {
					quit <- struct{}{}
				}
				return nil
			})
		go conn.Serve()
		<-quit
	case "client":
		dest := lib.DeserializeID(os.Args[2])
	}
}

func start(conn *lib.Connection) {
	conn.Write(data)
}

func initTime() {
	before = time.Now()
}

func calcElpasedTime(before time.Time) (time.Time, time.Duration) {
	now := time.Now()
	duration := now.Sub(before)
	return now, duration
}
