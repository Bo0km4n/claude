package lib

import (
	"encoding/binary"
	"errors"
	"log"
	"net"

	peerConfig "github.com/Bo0km4n/claude/app/peer/config"

	"github.com/Bo0km4n/claude/app/peer/service"
)

type Connection struct {
	NetConn           net.Conn
	Protocol          string
	DestinationPeerID []byte
	SourcePeerID      []byte
}

type ClaudePacket struct {
	DestinationPeerID [36]byte
	SourcePeerID      [36]byte
	CheckSum          uint16
	Payload           []byte
}

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
			SourcePeerID:      service.GetPeerID(),
		}, nil
	case "tcp":
		return &Connection{
			NetConn:           service.NetConn,
			Protocol:          "tcp",
			DestinationPeerID: dest,
			SourcePeerID:      service.GetPeerID(),
		}, nil
	}
	return nil, errors.New("Not found network")
}

func (c *Connection) Ping() {
	msg := []byte(`Hello world!`)
	packets := buildPacket(c, msg)
	if len(packets) == 0 {
		log.Printf("Packet is empty")
		return
	}
	buf := make([]byte, 1024)
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

func (cp *ClaudePacket) Serialize() []byte {
	b := append(cp.SourcePeerID[:], cp.DestinationPeerID[:]...)
	checkSumBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(checkSumBytes, cp.CheckSum)
	b = append(b, checkSumBytes...)
	b = append(b, cp.Payload...)
	return b
}
