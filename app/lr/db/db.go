package db

import (
	"fmt"

	"github.com/k0kubun/pp"

	"github.com/Bo0km4n/claude/app/common/proto"
)

var memcache map[string]*proto.PeerEntry

func InitDB() {
	memcache = make(map[string]*proto.PeerEntry, 1024)
}

func InsertEntry(key string, value *proto.PeerEntry) {
	memcache[key] = value
}

func FetchEntry(key string) (*proto.PeerEntry, error) {
	v, ok := memcache[key]
	if !ok {
		return nil, fmt.Errorf("Not found key: %s", key)
	}
	return v, nil
}

func Dump() {
	pp.Println(memcache)
}
