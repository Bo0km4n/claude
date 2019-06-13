package packet

const (
	PACKET_SIZE   = 0xFFFF // this size includes header length
	HEADER_LENGTH = 0x4c   // Decimal = 76
)

type ClaudePacket struct {
	DestinationPeerID [36]byte
	SourcePeerID      [36]byte
	CheckSum          uint16
	PayloadLen        uint16
	Payload           []byte
}
