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
	s := sha256.Sum256([]byte("abcd"))
	d := sha256.Sum256([]byte("efgh"))
	in := &Connection{
		SourcePeerID:      s[:],
		DestinationPeerID: d[:],
	}

	header, err := buildHeader(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(header) != 66 {
		t.Errorf("expected len(header) is 66, but got = %d", len(header))
	}
}

func TestWrapHeader(t *testing.T) {
	in := [][]byte{
		[]byte(`Hello world!`),
	}
	s := sha256.Sum256([]byte("abcd"))
	d := sha256.Sum256([]byte("efgh"))
	conn := &Connection{
		SourcePeerID:      s[:],
		DestinationPeerID: d[:],
	}
	packets, err := addHeader(conn, in)
	if err != nil {
		t.Fatal(err)
	}
	if len(packets) != 1 {
		t.Errorf("expected len(packets) = 1, but got = %d", len(packets))
	}
	if len(packets[0]) != 78 {
		t.Errorf("expected len(packets[0]) = 78, but got = %d", len(packets[0]))
	}
}
