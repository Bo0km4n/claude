package db

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/k0kubun/pp"
)

var memcache sync.Map

func InitDB() {
	memcache = sync.Map{}
}

func InsertPeerEntry(key []byte, value *proto.PeerEntry) {
	// keyStr := hex.EncodeToString(key)
	keyStr := base64.StdEncoding.EncodeToString(key)
	memcache.Store(keyStr, value)
}

func FetchPeerEntry(key []byte) (*proto.PeerEntry, error) {
	keyStr := base64.StdEncoding.EncodeToString(key)
	pp.Println(key)
	v, ok := memcache.Load(keyStr)
	if !ok {
		return nil, fmt.Errorf("Not found key: %s", key)
	}
	return v.(*proto.PeerEntry), nil
}

// PeerA: iNQmb9TmM40TuEX88olXnSCciXgjuSF9o+Fhk28DFYk=
// PeerB: 5eCIoLZhY6Cial4FPSpEltwWq24OPdGt8tFqqEoHjJ0=
func DebugInsertEntryPeerA() {
	key := "iNQmb9TmM40TuEX88olXnSCciXgjuSF9o+Fhk28DFYk="

	if _, ok := memcache.Load(key); ok {
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
	memcache.Store(key, v)
}

func DebugInsertEntryPeerB() {
	key := "5eCIoLZhY6Cial4FPSpEltwWq24OPdGt8tFqqEoHjJ0="
	if _, ok := memcache.Load(key); ok {
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
	memcache.Store(key, v)
}
