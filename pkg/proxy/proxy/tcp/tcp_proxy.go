package tcp

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/Bo0km4n/claude/claude/golang/packet"
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
	peerID, ok := pipe.FetchIdByIp(net.ParseIP(peerAddrStr).String())
	if !ok {
		log.Printf("Not found ip: %s", net.ParseIP(peerAddrStr).String())
		return
	}

	// New pipe connection
	p := &pipe.Pipe{
		Addr:               peerAddrStr,
		PeerConnection:     in,
		ProxyConnectionMap: map[string]net.Conn{},
	}
	pipe.Insert(peerID, p)
	if err := tp.upRelay(p); err != nil {
		log.Println(err)
		return
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
		log.Println("upRelay read | ", n)
		packets, err := packet.Parse(buf[:n])
		if err != nil {
			return err
		}
		for _, pac := range packets {
			destID := pac.GetDestinationID()
			proxyConn, ok := p.ProxyConnectionMap[destID]
			if !ok {
				newProxyConn, err := newConnectionToProxy(destID)
				if err != nil {
					return err
				}
				p.ProxyConnectionMap[destID] = newProxyConn
				proxyConn = newProxyConn
			}
			if _, err := proxyConn.Write(pac.Serialize()); err != nil {
				return err
			}
		}
	}
}

func (tp *TCPProxy) downHandleConn(in *net.TCPConn) {
	if err := tp.downRelay(in); err != nil {
		log.Fatal(err)
	}
}

func (tp *TCPProxy) downRelay(in *net.TCPConn) error {
	buf := make([]byte, packet.PACKET_SIZE)
	for {
		n, err := in.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("downRelay read | ", n)
		packets, err := packet.Parse(buf[:n])
		if err != nil {
			return err
		}
		for _, p := range packets {
			tp.relayToPeer(p)
		}
	}
}

func (tp *TCPProxy) relayToPeer(p *packet.ClaudePacket) {
	id := p.GetDestinationID()
	pipe, ok := pipe.Fetch(id)
	if !ok {
		log.Printf("%+v", xerrors.Errorf("Not found peer: %s", id))
		return
	}
	if _, err := pipe.PeerConnection.Write(p.Serialize()); err != nil {
		log.Printf("%+v", xerrors.Errorf("%+v", err))
		return
	}
}

func (tp *TCPProxy) serveUpStream() {
	listener, err := net.Listen("tcp", ":"+config.Config.Claude.UpTcpPort)
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
	listener, err := net.Listen("tcp", ":"+config.Config.Claude.DownTcpPort)
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
