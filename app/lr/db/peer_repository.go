package db

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"sync"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/k0kubun/pp"
)

var peerRepository sync.Map

func InitDB() {
	peerRepository = sync.Map{}
}

func InsertPeerEntry(key []byte, value *proto.PeerEntry) {
	keyStr := base64.StdEncoding.EncodeToString(key)
	peerRepository.Store(keyStr, value)
}

func FetchPeerEntry(key []byte) (*proto.PeerEntry, error) {
	keyStr := base64.StdEncoding.EncodeToString(key)
	v, ok := peerRepository.Load(keyStr)
	if !ok {
		return fetchPeerEntryFromTablet(binary.BigEndian.Uint32(key[0:4]))
	}
	return v.(*proto.PeerEntry), nil
}

func fetchPeerEntryFromTablet(id uint32) (*proto.PeerEntry, error) {
	pp.Println(id)
	return nil, errors.New("Not foune key in tablet server")
}

// PeerA: iNQmb9TmM40TuEX88olXnSCciXgjuSF9o+Fhk28DFYk=
// PeerB: 5eCIoLZhY6Cial4FPSpEltwWq24OPdGt8tFqqEoHjJ0=
func DebugInsertEntryPeerA() {
	key := "iNQmb9TmM40TuEX88olXnSCciXgjuSF9o+Fhk28DFYk="

	if _, ok := peerRepository.Load(key); ok {
		pp.Println("Skip debug insert peerA")
		return
	}

	keyB, _ := base64.StdEncoding.DecodeString(key)
	v := &proto.PeerEntry{
		PeerId:    keyB,
		LocalIp:   "100.100.100.100",
		LocalPort: "9610",
		IsRemote:  true,
	}
	peerRepository.Store(key, v)
}

func DebugInsertEntryPeerB() {
	key := "5eCIoLZhY6Cial4FPSpEltwWq24OPdGt8tFqqEoHjJ0="
	if _, ok := peerRepository.Load(key); ok {
		pp.Println("Skip debug insert peerB")
		return
	}

	keyB, _ := base64.StdEncoding.DecodeString(key)
	v := &proto.PeerEntry{
		PeerId:    keyB,
		LocalIp:   "100.100.100.200",
		LocalPort: "9610",
		IsRemote:  true,
	}
	peerRepository.Store(key, v)
}
