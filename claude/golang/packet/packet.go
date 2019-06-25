package packet

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	PACKET_SIZE   = 0xffff // this size includes header length
	HEADER_LENGTH = 0x4e   // Decimal = 78
	ID_LENGTH     = 36
)

const (
	CONTROL_FLAG_PING = uint16(iota)
	CONTROL_FLAG_NORMAL
)

type ClaudePacket struct {
	header  *claudeHeader
	payload payload
}

type claudeHeader struct {
	ControlFlag       uint16
	DestinationPeerID [ID_LENGTH]byte
	SourcePeerID      [ID_LENGTH]byte
	CheckSum          uint16
	PayloadLen        uint16
}

type payload []byte

func Parse(b []byte) ([]*ClaudePacket, error) {
	if len(b) < HEADER_LENGTH {
		return nil, fmt.Errorf("Header length is too short: %d", len(b))
	}
	packets := []*ClaudePacket{}
	r := bytes.NewReader(b)
	for {
		hbuf := make([]byte, HEADER_LENGTH)
		_, err := r.Read(hbuf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			return packets, nil
		}
		h := ParseHeader(hbuf)
		payload := make([]byte, h.PayloadLen)
		_, err = r.Read(payload)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			return packets, nil
		}
		p := &ClaudePacket{
			header:  h,
			payload: payload,
		}
		packets = append(packets, p)
	}
	return packets, nil
}

func ParseHeader(header []byte) *claudeHeader {
	h := &claudeHeader{
		ControlFlag: binary.BigEndian.Uint16(header[0:2]),
		CheckSum:    binary.BigEndian.Uint16(header[74:76]),
		PayloadLen:  binary.BigEndian.Uint16(header[76:78]),
	}
	copy(h.DestinationPeerID[:], header[2:38])
	copy(h.SourcePeerID[:], header[38:74])
	return h
}

func GeneratePingPacket() *ClaudePacket {
	return &ClaudePacket{
		header: &claudeHeader{
			ControlFlag:       CONTROL_FLAG_PING,
			DestinationPeerID: [ID_LENGTH]byte{},
			SourcePeerID:      [ID_LENGTH]byte{},
			CheckSum:          0,
			PayloadLen:        0,
		},
		payload: []byte{},
	}
}

func GeneratePacket() *ClaudePacket {
	return &ClaudePacket{
		header: &claudeHeader{
			ControlFlag:       CONTROL_FLAG_NORMAL,
			DestinationPeerID: [ID_LENGTH]byte{},
			SourcePeerID:      [ID_LENGTH]byte{},
			CheckSum:          0,
			PayloadLen:        0,
		},
		payload: []byte{},
	}
}

func (cp *ClaudePacket) SetToID(id string) error {
	idBytes, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return err
	}
	if len(idBytes) != ID_LENGTH {
		return fmt.Errorf("Not matched id length=%d", len(idBytes))
	}
	copy(cp.header.DestinationPeerID[:], idBytes[:])
	return nil
}

func (cp *ClaudePacket) SetFromID(id string) error {
	idBytes, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return err
	}
	if len(idBytes) != ID_LENGTH {
		return fmt.Errorf("Not matched id length=%d", len(idBytes))
	}
	copy(cp.header.SourcePeerID[:], idBytes[:])
	return nil
}

func (cp *ClaudePacket) SetPayload(p []byte) error {
	if len(p) > PACKET_SIZE-HEADER_LENGTH {
		return fmt.Errorf("Too large payload size=%d", len(p))
	}
	cp.header.PayloadLen = uint16(len(p))
	cp.payload = p
	return nil
}

func (cp *ClaudePacket) SetHeader(h *claudeHeader) {
	cp.header = h
}

func (cp *ClaudePacket) GetDestinationID() string {
	return base64.StdEncoding.EncodeToString(cp.header.DestinationPeerID[:])
}

func (cp *ClaudePacket) Serialize() []byte {
	b := make([]byte, 0)
	u16Bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(u16Bytes, cp.header.ControlFlag)
	b = append(b, u16Bytes...)
	b = append(b, cp.header.DestinationPeerID[:]...)
	b = append(b, cp.header.SourcePeerID[:]...)
	binary.BigEndian.PutUint16(u16Bytes, cp.header.CheckSum)
	b = append(b, u16Bytes...)
	binary.BigEndian.PutUint16(u16Bytes, cp.header.PayloadLen)
	b = append(b, u16Bytes...)
	b = append(b, cp.payload...)
	return b
}

func (cp *ClaudePacket) PayloadLen() uint16 {
	return cp.header.PayloadLen
}

type PacketBuffer struct {
	Buffer chan []byte
}
