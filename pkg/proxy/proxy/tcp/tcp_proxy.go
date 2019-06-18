package tcp

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/Bo0km4n/claude/claude/go/packet"
	"github.com/Bo0km4n/claude/pkg/proxy/config"
	"github.com/Bo0km4n/claude/pkg/proxy/repository/pipe"
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
	p := &pipe.Pipe{
		Addr:           peerAddrStr,
		PeerConnection: in,
	}
	pipe.Insert(peerID, p)
	if err := tp.upRelay(p); err != nil {
		log.Fatal(err)
	}
	return
}

func (tp *TCPProxy) upRelay(p *pipe.Pipe) error {
	buf := make([]byte, packet.PACKET_SIZE)
	defer p.PeerConnection.Close()
	for {
		n, err := p.PeerConnection.Read(buf)
		if err != nil {
			log.Println(n, err)
			return err
		}
		log.Println("Read: ", n)
		b := buf[:n]

		if p.ProxyConnection == nil {
			// Maybe when first read, proxy connection is not established yet.
			// So connect to remote proxy and store pipe

			// TODO: set parsed id
			proxyConn, err := newConnectionToProxy("hoge")
			if err != nil {
				return err
			}
			p.ProxyConnection = proxyConn
		}
		n, err = p.ProxyConnection.Write(b)
		if err != nil {
			return err
		}
		log.Println("Write: ", n)
	}
}

func (tp *TCPProxy) downHandleConn(in *net.TCPConn) {

}

func (tp *TCPProxy) downRelay(in *net.TCPConn) error {
	buf := make([]byte, packet.PACKET_SIZE)
	for {
		n, err := in.Read(buf)
		if err != nil {
			return err
		}
		log.Println("Read: ", n)
		// b := buf[:n]

		// TODO: Parse header to claude packet
		// packets := parseClaudePackets(b)
		// for _, p := range {
		// 	tp.relayToPeer(p)
		// }
	}
}

func (tp *TCPProxy) relayToPeer(p *packet.ClaudePacket) {
	// id := p.GetDestinationPeerID()
	// pipe, ok := repository.FetchPipe(id)
	// if !ok {
	// 	log.Printf("Not found peer: %s\n", id)
	// 	return
	// }
	// pipe.PeerConnection.Write(p.Serialize())
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
