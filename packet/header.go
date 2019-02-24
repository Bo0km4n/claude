package packet

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Header is packet header defined by claude and contains data
type Header struct {
	layers.BaseLayer

	SrcIP    uint32
	DstIP    uint32
	SrcPort  uint16
	DstPort  uint16
	CheckSum uint16
	Length   uint16
}

var (
	// LayerTypeClaude implement gopacket layer type
	LayerTypeClaude = gopacket.RegisterLayerType(200, gopacket.LayerTypeMetadata{Name: "CLAUDE", Decoder: gopacket.DecodeFunc(decodeClaudeHeader)})
)

// LayerType returns header.ClaudeHeader
func (h *Header) LayerType() gopacket.LayerType { return LayerTypeClaude }

// decodeClaudeHeader decodes the byte slice into a Claude Header type.
// It also setups the application Layer in PacketBuilder
func decodeClaudeHeader(data []byte, p gopacket.PacketBuilder) error {
	h := &Header{}
	err := h.DecodeFromBytes(data, p)
	if err != nil {
		return err
	}
	p.AddLayer(h)
	p.SetApplicationLayer(h)
	return nil
}

// DecodeFromBytes decodes the slice into the Claude Header Struct
func (h *Header) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 16 {
		df.SetTruncated()
		return errors.New("Claude packet too short")
	}

	h.BaseLayer = layers.BaseLayer{Contents: data[:16]}
	h.SrcIP = binary.BigEndian.Uint32(data[:4])
	h.DstIP = binary.BigEndian.Uint32(data[4:8])
	h.SrcPort = binary.BigEndian.Uint16(data[8:10])
	h.DstPort = binary.BigEndian.Uint16(data[10:12])
	h.CheckSum = binary.BigEndian.Uint16(data[12:14])
	h.Length = binary.BigEndian.Uint16(data[14:16])

	switch {
	case h.Length >= 16:
		hlen := int(h.Length)
		if hlen > len(data) {
			df.SetTruncated()
			hlen = len(data)
		}
		h.BaseLayer.Payload = data[16:hlen]
	case h.Length == 0:
		h.BaseLayer.Payload = data[16:]
	default:
		return fmt.Errorf("CLAUDE packet too small: %d bytes", h.Length)
	}
	return nil
}

// CanDecode implements gopacket.DecodingLayer.
func (h *Header) CanDecode() gopacket.LayerClass {
	return LayerTypeClaude
}

// Payload returns Header.BaseLayer.Payload
func (h *Header) Payload() []byte {
	return h.BaseLayer.Payload
}

// NextLayerType implements gopacket.DecodingLayer.
func (h *Header) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypePayload
}
