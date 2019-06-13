package proxy

import (
	"log"
	"net"
	"sync"

	"github.com/Bo0km4n/claude/pkg/proxy/config"
	"golang.org/x/xerrors"
)

type TCPProxy struct {
	wg sync.WaitGroup
}

func (tp *TCPProxy) upHandleConn(in *net.TCPConn) {
	// Register peer information
	// Set peer handler
}
func (tp *TCPProxy) downHandleConn(in *net.TCPConn) {

}

func (tp *TCPProxy) serveUpStream() {
	listener, err := net.Listen("tcp", "localhost:"+config.Config.Claude.UpTcpPort)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(xerrors.Errorf("%+v", err))
		}
		tcpConn, _ := conn.(*net.TCPConn)
		go tp.upHandleConn(tcpConn)
	}
}

func (tp *TCPProxy) serveDownStream() {
	listener, err := net.Listen("tcp", "localhost:"+config.Config.Claude.DownTcpPort)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(xerrors.Errorf("%+v", err))
		}
		tcpConn, _ := conn.(*net.TCPConn)
		go tp.downHandleConn(tcpConn)
	}
}

func (tp *TCPProxy) Serve() {
	tp.serveUpStream()
	tp.wg.Add(1)
	tp.serveDownStream()
	tp.wg.Add(1)

	tp.wg.Wait()
}
