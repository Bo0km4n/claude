package db

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/Bo0km4n/claude/app/common/proto"
)

var memcache sync.Map

func InitDB() {
	memcache = sync.Map{}
}

func InsertEntry(key []byte, value *proto.PeerEntry) {
	keyStr := hex.EncodeToString(key)
	memcache.Store(keyStr, value)
}

func FetchEntry(key []byte) (*proto.PeerEntry, error) {
	keyStr := hex.EncodeToString(key)
	v, ok := memcache.Load(keyStr)
	if !ok {
		return nil, fmt.Errorf("Not found key: %s", key)
	}
	return v.(*proto.PeerEntry), nil
}
