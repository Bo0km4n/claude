package packet

import "fmt"

const (
	PACKET_SIZE   = 0xFFFF // this size includes header length
	HEADER_LENGTH = 0x4e   // Decimal = 78
)

const (
	CONTROL_FLAG_PING = uint16(iota)
	CONTROL_FLAG_NORMAL
)

type ClaudePacket struct {
	ControlFlag       uint16
	DestinationPeerID [36]byte
	SourcePeerID      [36]byte
	CheckSum          uint16
	PayloadLen        uint16
	Payload           []byte
}

func Parse(b []byte) (*ClaudePacket, error) {
	if len(b) < HEADER_LENGTH {
		return nil, fmt.Errorf("Header length is too short: %d", len(b))
	}
	return nil, nil
}

func GeneratePingPacket() *ClaudePacket {
	return &ClaudePacket{
		ControlFlag:       CONTROL_FLAG_PING,
		DestinationPeerID: [36]byte{},
		SourcePeerID:      [36]byte{},
		CheckSum:          0,
		PayloadLen:        0,
		Payload:           []byte{},
	}
}
