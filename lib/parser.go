package lib

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func ParseHeader(payload []byte) (*ClaudePacket, error) {
	if len(payload) < 66 {
		return nil, fmt.Errorf("Packet's length is too short, got = %d", len(payload))
	}
	srcID := payload[0:32]
	dstID := payload[32:64]
	checkSum := binary.BigEndian.Uint16(payload[64:66])

	if !validCheckSum(checkSum, srcID, dstID) {
		return nil, errors.New("Invalid checksum")
	}

	cp := &ClaudePacket{
		CheckSum: checkSum,
	}
	copy(cp.SourcePeerID[:], srcID)
	copy(cp.DestinationPeerID[:], dstID)
	cp.Payload = payload[66:]
	return cp, nil
}
