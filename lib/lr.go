package lib

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/Bo0km4n/claude/app/peer/service"
)

func ConnectToLR(protocol string) {
	done := make(chan int)
	go service.LaunchGRPCService(done, protocol)
	<-done

	time.Sleep(2)
	service.UDPBcast()
	for {
		if service.IsCompletedJoinToLR {
			return
		}
		time.Sleep(1)
	}
}

func CryptedID(seed string) []byte {
	seed256 := sha256.Sum256([]byte(seed))
	lrID := make([]byte, 4)
	binary.BigEndian.PutUint32(lrID, service.RemoteLR.ID)
	dest := append(lrID, seed256[:]...)
	return dest
}

func DeserializeID(id string) []byte {
	b, _ := base64.StdEncoding.DecodeString(id)
	return b
}
