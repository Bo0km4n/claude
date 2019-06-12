package lib

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/Bo0km4n/claude/app/common/proto"

	"github.com/Bo0km4n/claude/app/peer/api"
	peerConfig "github.com/Bo0km4n/claude/app/peer/config"

	"github.com/Bo0km4n/claude/app/peer/service"
)

type Connection struct {
	NetConn           net.Conn
	Protocol          string
	DestinationPeerID []byte
	SourcePeerID      []byte
	Handler           func(*Connection, []byte) error
}

type ClaudePacket struct {
	DestinationPeerID [36]byte
	SourcePeerID      [36]byte
	CheckSum          uint16
	Payload           []byte
}

const BUF_SIZE = 1024

func InitConfig() {
	peerConfig.InitConfig()
}

func NewConnection(protocol string, dest []byte) (*Connection, error) {

	switch protocol {
	case "udp":
		return &Connection{
			NetConn:           service.NetConn,
			Protocol:          "udp",
			DestinationPeerID: dest,
			SourcePeerID:      api.GetPeerIDBytes(),
		}, nil
	case "tcp":
		return &Connection{
			NetConn:           service.NetConn,
			Protocol:          "tcp",
			DestinationPeerID: dest,
			SourcePeerID:      api.GetPeerIDBytes(),
		}, nil
	}

	return nil, errors.New("Not found network")
}

func (c *Connection) LookUp(distance float32) []*proto.PeerEntry {
	return nil
}

func (c *Connection) Ping() {
	msg := []byte(`Hello world!`)
	packets := buildPacket(c, msg)
	if len(packets) == 0 {
		log.Printf("Packet is empty")
		return
	}
	buf := make([]byte, BUF_SIZE)
	for _, p := range packets {
		n, err := c.NetConn.Write(p)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Send packets: %d\n", n)
		}
		for {
			n, err = c.NetConn.Read(buf)
			if err != nil {
				log.Fatal(err)
			} else {
				resp, err := ParseHeader(buf)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Received msg: %s", string(resp.Payload))
			}
		}
	}
}

func (c *Connection) RegisterHandler(f func(*Connection, []byte) error) {
	c.Handler = f
}

func (c *Connection) Serve() error {
	c.NetConn.Write([]byte{0x00})
	for {
		buf := make([]byte, BUF_SIZE)
		_, err := c.NetConn.Read(buf)
		if err != nil {
			return err
		} else {
			resp, err := ParseHeader(buf)
			if err != nil {
				return err
			}
			if err := c.Handler(c, resp.Payload); err != nil {
				return err
			}
		}
	}
}

func (c *Connection) Write(b []byte) error {
	packet := buildPacket(c, b)
	for i := range packet {
		if _, err := c.NetConn.Write(packet[i]); err != nil {
			return err
		}
	}
	return nil
}

type netConnFormat struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

func (c *Connection) SaveConnection() {
	laddr := c.NetConn.LocalAddr()
	switch c.Protocol {
	case "tcp":
		tcpAddr := laddr.(*net.TCPAddr)
		f := &netConnFormat{
			IP:       tcpAddr.IP.String(),
			Port:     fmt.Sprintf("%d", tcpAddr.Port),
			Protocol: c.Protocol,
		}
		fb, err := json.Marshal(f)
		if err != nil {
			log.Fatal(err)
		}
		c.dumpConnectionToFile(fb)
	case "udp":
		udpAddr := laddr.(*net.UDPAddr)
		f := &netConnFormat{
			IP:       udpAddr.IP.String(),
			Port:     fmt.Sprintf("%d", udpAddr.Port),
			Protocol: c.Protocol,
		}
		fb, err := json.Marshal(f)
		if err != nil {
			log.Fatal(err)
		}
		c.dumpConnectionToFile(fb)
	}
}

func (c *Connection) dumpConnectionToFile(b []byte) error {
	file, err := os.OpenFile(peerConfig.Config.NetConnFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(b)
	return err
}

func (c *Connection) RestoreConnection() error {
	ncf, err := c.loadFile()
	if err != nil {
		return err
	}
	ip := net.ParseIP(ncf.IP)
	port, err := strconv.Atoi(ncf.Port)
	if err != nil {
		return err
	}
	switch ncf.Protocol {
	case "tcp":
		laddr := &net.TCPAddr{
			IP:   ip,
			Port: port,
		}
		raddr, err := net.ResolveTCPAddr("tcp", service.RemoteLR.Addr+":"+service.RemoteLR.TcpPort)
		if err != nil {
			return err
		}
		conn, err := net.DialTCP("tcp", laddr, raddr)
		if err != nil {
			return err
		}
		c.NetConn = conn
	case "udp":
		laddr := &net.UDPAddr{
			IP:   ip,
			Port: port,
		}
		raddr, err := net.ResolveUDPAddr("udp", service.RemoteLR.Addr+":"+service.RemoteLR.TcpPort)
		if err != nil {
			return err
		}
		conn, err := net.DialUDP("udp", laddr, raddr)
		if err != nil {
			return err
		}
		c.NetConn = conn
	}

	return errors.New("Unexpected protocol")
}

func (c *Connection) loadFile() (*netConnFormat, error) {
	file, err := os.Open(peerConfig.Config.NetConnFile)
	if err != nil {
		return &netConnFormat{}, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return &netConnFormat{}, err
	}
	ncf := &netConnFormat{}
	if err := json.Unmarshal(b, ncf); err != nil {
		return ncf, err
	}
	return ncf, nil
}

func (cp *ClaudePacket) Serialize() []byte {
	b := append(cp.SourcePeerID[:], cp.DestinationPeerID[:]...)
	checkSumBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(checkSumBytes, cp.CheckSum)
	b = append(b, checkSumBytes...)
	b = append(b, cp.Payload...)
	return b
}
