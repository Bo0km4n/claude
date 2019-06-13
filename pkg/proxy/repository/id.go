package repository

import (
	"net"
	"sync"
)

type IDRepo struct {
	Map map[string]*Pipe // Key: ID, Value: Pipe connection
	mu  sync.Mutex
}

type Pipe struct {
	Addr            string // IP:Port
	PeerConnection  net.Conn
	ProxyConnection net.Conn
}

var idRepo *IDRepo

func InsertPipe(key string, value *Pipe) {
	idRepo.mu.Lock()
	defer idRepo.mu.Unlock()
	idRepo.Map[key] = value
}

func FetchPipe(key string) (*Pipe, bool) {
	idRepo.mu.Lock()
	defer idRepo.mu.Unlock()
	v, ok := idRepo.Map[key]
	return v, ok
}
