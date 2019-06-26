package tcp

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Bo0km4n/claude/claude/golang/packet"
	"github.com/Bo0km4n/claude/pkg/proxy/config"
	"github.com/Bo0km4n/claude/pkg/proxy/repository/pipe"
	"golang.org/x/xerrors"
)

var idleTimeout = time.Second * 2400

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
	defer p.PeerConnection.Close()
	for {
		headerBuf := make([]byte, packet.HEADER_LENGTH)
		if _, err := io.ReadFull(p.PeerConnection, headerBuf); err != nil {
			return err
		}
		header := packet.ParseHeader(headerBuf)
		payloadBuf := make([]byte, int(header.PayloadLen))
		if _, err := io.ReadFull(p.PeerConnection, payloadBuf); err != nil {
			return err
		}
		newPacket := packet.GeneratePacket()
		newPacket.SetPayload(payloadBuf)
		newPacket.SetHeader(header)
		destID := newPacket.GetDestinationID()
		proxyConn, ok := p.ProxyConnectionMap[destID]
		if !ok {
			newProxyConn, err := newConnectionToProxy(destID)
			if err != nil {
				return err
			}
			p.ProxyConnectionMap[destID] = newProxyConn
			proxyConn = newProxyConn
		}
		n, err := proxyConn.Write(newPacket.Serialize())
		if err != nil {
			return err
		}
		log.Printf("upRelay write: %d\n", n)
	}
}

func (tp *TCPProxy) downHandleConn(in *net.TCPConn) {
	err := tp.downRelay(in)
	if err == io.EOF {
		in.Close()
		return
	}
	if err != nil {
		log.Fatal("downHandleConn: ", err)
	}
}

func (tp *TCPProxy) downRelay(in *net.TCPConn) error {
	for {
		headerBuf := make([]byte, packet.HEADER_LENGTH)
		if _, err := io.ReadFull(in, headerBuf); err != nil {
			return err
		}
		header := packet.ParseHeader(headerBuf)
		payloadBuf := make([]byte, int(header.PayloadLen))
		if _, err := io.ReadFull(in, payloadBuf); err != nil {
			return err
		}
		newPacket := packet.GeneratePacket()
		newPacket.SetPayload(payloadBuf)
		newPacket.SetHeader(header)
		tp.relayToPeer(newPacket)
	}
}

func (tp *TCPProxy) relayToPeer(p *packet.ClaudePacket) {
	id := p.GetDestinationID()
	pipe, ok := pipe.Fetch(id)
	if !ok {
		log.Printf("%+v", xerrors.Errorf("Not found peer: %s", id))
		return
	}
	n, err := pipe.PeerConnection.Write(p.Serialize())
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("%+v", err))
		return
	}
	log.Printf("downRelay write: %d\n", n)
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
		tcpConn.SetDeadline(time.Now().Add(idleTimeout))
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
		tcpConn.SetDeadline(time.Now().Add(idleTimeout))
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
