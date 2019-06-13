package tcp

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
	defer in.Close()
	peerAddrStr := in.RemoteAddr().String()
	peerID, err := getPeerID(net.ParseIP(peerAddrStr).String())
	if err != nil {
		log.Println(err)
		return
	}

	// New pipe connection
	pipe := &repository.Pipe{
		Addr:           peerAddrStr,
		PeerConnection: in,
	}
	repository.InsertPipe(peerID, pipe)
	if err := tp.upRelay(pipe); err != nil {
		log.Fatal(err)
	}
	return
}

func (tp *TCPProxy) upRelay(pipe *repository.Pipe) error {
	buf := make([]byte, packet.PACKET_SIZE)
	defer pipe.PeerConnection.Close()
	for {
		n, err := pipe.PeerConnection.Read(buf)
		if err != nil {
			log.Println(n, err)
			return err
		}
		log.Println("Read: ", n)
		b := buf[:n]

		if pipe.ProxyConnection == nil {
			// Maybe when first read, proxy connection is not established yet.
			// So connect to remote proxy and store pipe

			// TODO: set parsed id
			proxyConn, err := newConnectionToProxy("hoge")
			if err != nil {
				return err
			}
			pipe.ProxyConnection = proxyConn
		}
		n, err = pipe.ProxyConnection.Write(b)
		if err != nil {
			return err
		}
		log.Println("Write: ", n)
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

func NewProxy() *TCPProxy {
	return &TCPProxy{}
}
