package lib

import (
	"crypto/sha256"
	"testing"
)

func TestCalcCheckSum(t *testing.T) {
	in := []uint16{
		0x4500, 0x16ce, 0x654c, 0x4000,
		0x0111, 0xc0a8, 0x651f,
		0xe000, 0x001f,
	}

	result := calcCheckSum(in)
	if uint16(0xf7eb) != result {
		t.Errorf("expected %4x. But got %4x\n", 0xf7eb, result)
	}
}

func TestBuildHeader(t *testing.T) {
	srcPeerID := sha256.Sum256([]byte("abcd"))
	dstPeerID := sha256.Sum256([]byte("efgh"))
	srcProxyID := []byte{0x00, 0x00, 0x00, 0x01}
	dstProxyID := []byte{0x00, 0x00, 0x00, 0x02}

	in := &Connection{
		SourcePeerID:      append(srcProxyID, srcPeerID[:]...),
		DestinationPeerID: append(dstProxyID, dstPeerID[:]...),
	}

	header, err := buildHeader(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(header) != 74 {
		t.Errorf("expected len(header) is 74, but got = %d", len(header))
	}
}

func TestWrapHeader(t *testing.T) {
	in := [][]byte{
		[]byte(`Hello world!`),
	}
	srcPeerID := sha256.Sum256([]byte("abcd"))
	dstPeerID := sha256.Sum256([]byte("efgh"))
	srcProxyID := []byte{0x00, 0x00, 0x00, 0x01}
	dstProxyID := []byte{0x00, 0x00, 0x00, 0x02}
	conn := &Connection{
		SourcePeerID:      append(srcProxyID, srcPeerID[:]...),
		DestinationPeerID: append(dstProxyID, dstPeerID[:]...),
	}
	packets, err := addHeader(conn, in)
	if err != nil {
		t.Fatal(err)
	}
	if len(packets) != 1 {
		t.Errorf("expected len(packets) = 1, but got = %d", len(packets))
	}
	if len(packets[0]) != 86 {
		t.Errorf("expected len(packets[0]) = 86, but got = %d", len(packets[0]))
	}
}
