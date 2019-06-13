package proxy

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/Bo0km4n/claude/claude/go/packet"
	"github.com/Bo0km4n/claude/pkg/proxy/config"
	"github.com/Bo0km4n/claude/pkg/proxy/repository"
	"golang.org/x/xerrors"
)

type TCPProxy struct {
	wg sync.WaitGroup
}

func (tp *TCPProxy) upHandleConn(in *net.TCPConn) {
	// Register peer information
	// Set peer handler
	defer in.Close()
	peerAddr := in.RemoteAddr().String()
	idAndConn, ok := repository.FetchIDAndConn(peerAddr)

	if !ok {
		// Init peer connection and remote proxy connection
	}
	if idAndConn.ID != "" && idAndConn.RemoteProxyConn == nil {
		// Peer registered, but not yet connected remote proxy with dest id
	}
	if idAndConn.ID != "" && idAndConn.RemoteProxyConn != nil {
		// Relay packet from peer to remote proxy
		tcpRemoteConn, _ := idAndConn.RemoteProxyConn.(*net.TCPConn)
		if err := tp.upRelay(in, tcpRemoteConn); err != nil {
			log.Fatal(err)
		}
	}
}

func (tp *TCPProxy) upRelay(from, to *net.TCPConn) error {
	buff := make([]byte, packet.PACKET_SIZE)
	for {
		n, err := from.Read(buff)
		if err != nil {
			log.Println(n, err)
			return err
		}
		log.Println("Read: ", n)
		b := buff[:n]

		_, err = to.Write(b)
		if err != nil {
			return err
		}
	}
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
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go tp.serveUpStream()
	go tp.serveDownStream()
	<-quit
	log.Println("Interrupted Proxy Server")
}

func NewTCPProxy() *TCPProxy {
	return &TCPProxy{}
}
