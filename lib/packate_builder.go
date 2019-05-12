package lib

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	PACKET_CAPACITY = 1024
)

func buildPacket(conn *Connection, body []byte) [][]byte {
	plainPackets := make([][]byte, 0)
	if len(body) > PACKET_CAPACITY {
		// split packate to some packates.
		offset := PACKET_CAPACITY
		start := 0
		for {
			if offset >= len(body) {
				offset = len(body)
			}
			plainPackets = append(plainPackets, body[start:offset])
			if offset >= len(body) {
				break
			} else {
				start = offset
				offset = offset + PACKET_CAPACITY
			}
		}
	} else {
		plainPackets = append(plainPackets, body)
	}

	wrappedPackets, err := addHeader(conn, plainPackets)
	if err != nil {
		return [][]byte{}
	}

	return wrappedPackets
}

func addHeader(conn *Connection, packets [][]byte) ([][]byte, error) {
	wrappedPackets := [][]byte{}
	header, err := buildHeader(conn)
	if err != nil {
		return [][]byte{}, err
	}
	for i := range packets {
		p := append(header, packets[i]...)
		wrappedPackets = append(wrappedPackets, p)
	}
	return wrappedPackets, nil
}

// Build the claude header by big endian
func buildHeader(conn *Connection) ([]byte, error) {
	srcID := conn.SourcePeerID
	dstID := conn.DestinationPeerID

	if len(dstID) != 36 || len(srcID) != 36 {
		return []byte{}, fmt.Errorf("ID's length is not 36, got = %d, %d", len(dstID), len(srcID))
	}

	buff := append(srcID, dstID...)
	header16, err := convert8To16(buff)
	if err != nil {
		return []byte{}, err
	}

	checkSum := calcCheckSum(header16[:])

	checkSumBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(checkSumBytes, checkSum)

	header := append(srcID, dstID...)
	header = append(header, checkSumBytes...)
	return header, nil
}

func calcCheckSum(payload []uint16) uint16 {
	var checkSum uint32
	for i := range payload {
		checkSum += uint32(payload[i])
		if checkSum > 0x0000ffff {
			carry := checkSum >> 16
			checkSum = checkSum & 0x0000ffff
			checkSum += carry
		}
	}

	result := ^uint16(checkSum)
	return result
}

func convert8To16(header []byte) ([]uint16, error) {
	i := 0
	u16 := []uint16{}
	if len(header)%2 != 0 {
		return []uint16{}, errors.New("Can't convert bytes")
	}
	for {
		if i >= len(header) {
			break
		}
		u16 = append(u16, binary.BigEndian.Uint16([]byte{header[i], header[i+1]}))
		i = i + 2
	}
	return u16, nil
}

func validCheckSum(checkSum uint16, srcID, dstID []byte) bool {
	u16IDs, err := convert8To16(append(srcID, dstID...))
	if err != nil {
		return false
	}
	return checkSum == calcCheckSum(u16IDs)
}
