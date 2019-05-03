package db

import (
	"fmt"
	"sync"

	"github.com/Bo0km4n/claude/app/common/proto"
)

var memcache sync.Map

func InitDB() {
	memcache = sync.Map{}
}

func InsertEntry(key []byte, value *proto.PeerEntry) {
	memcache.Store(key, value)
}

func FetchEntry(key []byte) (*proto.PeerEntry, error) {
	v, ok := memcache.Load(key)
	if !ok {
		return nil, fmt.Errorf("Not found key: %s", key)
	}
	return v.(*proto.PeerEntry), nil
}
