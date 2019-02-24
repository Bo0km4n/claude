package packet

import (
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var testClaudePacket = []byte{
	0x08, 0x00, 0x27, 0x90, 0x33, 0x4c, 0x08, 0x00, 0x27, 0xc1, 0x7c, 0x8c, 0x08, 0x00, 0x45, 0x00,
	0x00, 0x37, 0x0a, 0xf4, 0x40, 0x00, 0x40, 0x11, 0x99, 0xa8, 0xc0, 0xa8, 0x0a, 0x65, 0xc0, 0xa8,
	0x0a, 0x64, 0xd4, 0xd8, 0xc3, 0x50, 0x00, 0x23, 0x6d, 0xab, 0xc8, // eth, ip, udp headers

	// start claude header and message
	0x0a, 0x01, 0x01, 0xa0, 0x0a, 0x01, 0x02, 0xc3, 0x5a, 0xc3, 0x5b, 0x01, 0x01, 0x00, 0x1b,
	// payload = "Hello world"
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
}

func TestClaudePacket(t *testing.T) {
	p := gopacket.NewPacket(testClaudePacket, layers.LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{layers.LayerTypeEthernet, layers.LayerTypeIPv4, layers.LayerTypeUDP}, t)
}

func checkLayers(p gopacket.Packet, want []gopacket.LayerType, t *testing.T) {
	layers := p.Layers()
	t.Log("Checking packet layers, want", want)
	for _, l := range layers {
		t.Logf("  Got layer %v, %d bytes, payload of %d bytes", l.LayerType(),
			len(l.LayerContents()), len(l.LayerPayload()))
	}
	t.Log(p)
	if len(layers) < len(want) {
		t.Errorf("  Number of layers mismatch: got %d want %d", len(layers),
			len(want))
		return
	}
	for i, l := range want {
		if l == gopacket.LayerTypePayload {
			// done matching layers
			return
		}

		if layers[i].LayerType() != l {
			t.Errorf("  Layer %d mismatch: got %v want %v", i,
				layers[i].LayerType(), l)
		}
	}
}
