package repository

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"sync"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/k0kubun/pp"
)

type peerBucket struct {
	mu sync.Mutex
	b  map[string]*proto.PeerEntry
}

func (pb *peerBucket) Store(key string, v *proto.PeerEntry) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.b[key] = v
}

func (pb *peerBucket) Load(key string) (*proto.PeerEntry, bool) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	v, ok := pb.b[key]
	return v, ok
}

func (pb *peerBucket) Values() []*proto.PeerEntry {
	values := []*proto.PeerEntry{}
	for _, v := range pb.b {
		values = append(values, v)
	}
	return values
}

func (pb *peerBucket) Dump() {
	pp.Println(pb.b)
}

var peerRepository peerBucket

func InitDB() {
	peerRepository = peerBucket{
		mu: sync.Mutex{},
		b:  make(map[string]*proto.PeerEntry, 1024),
	}
}

func InsertPeerEntry(key []byte, value *proto.PeerEntry) {
	keyStr := base64.StdEncoding.EncodeToString(key)
	pp.Println(">>> new peer:", keyStr, value.LocalPort)
	peerRepository.Store(keyStr, value)
}

func FetchPeerEntry(key []byte) (*proto.PeerEntry, error) {
	keyStr := base64.StdEncoding.EncodeToString(key)
	v, ok := peerRepository.Load(keyStr)

	if !ok {
		pp.Println("Not found id:", keyStr)
		return fetchPeerEntryFromTablet(binary.BigEndian.Uint32(key[0:4]))
	}
	return v, nil
}

func FetchLocalPeers() []*proto.PeerEntry {
	peers := []*proto.PeerEntry{}
	values := peerRepository.Values()
	for _, v := range values {
		if !v.IsRemote {
			peers = append(peers, v)
		}
	}

	return peers
}

func Dump() {
	peerRepository.Dump()
}

func fetchPeerEntryFromTablet(id uint32) (*proto.PeerEntry, error) {
	return nil, errors.New("Not found key in tablet server")
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
