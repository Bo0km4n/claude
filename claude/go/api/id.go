package api

import (
	"encoding/base64"

	"github.com/Bo0km4n/claude/claude/go/service"
)

func GetPeerIDString() string {
	return service.PeerSvc.ID
}

func GetPeerIDBytes() []byte {
	b, _ := base64.StdEncoding.DecodeString(service.PeerSvc.ID)
	return b
}
