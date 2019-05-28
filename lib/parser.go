package lib

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func ParseHeader(payload []byte) (*ClaudePacket, error) {
	if len(payload) < 74 {
		return nil, fmt.Errorf("Packet's length is too short, got = %d", len(payload))
	}
	srcID := payload[0:36]
	dstID := payload[36:72]
	checkSum := binary.BigEndian.Uint16(payload[72:74])

	if !validCheckSum(checkSum, srcID, dstID) {
		return nil, errors.New("Invalid checksum")
	}

	cp := &ClaudePacket{
		CheckSum: checkSum,
	}
	copy(cp.SourcePeerID[:], srcID)
	copy(cp.DestinationPeerID[:], dstID)
	cp.Payload = payload[74:]
	return cp, nil
}
