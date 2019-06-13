package repository

import (
	"net"
	"sync"
)

type IPPortRepo struct {
	Map map[string]*IDAndConn // Key: IP:Port, Value: ID and Connection
	mu  sync.Mutex
}

type IDAndConn struct {
	ID             string
	FromConnection net.Conn
	ToConnection   net.Conn
}

var ipPortRepository *IPPortRepo

func FetchIDAndConn(key string) (*IDAndConn, bool) {
	ipPortRepository.mu.Lock()
	defer ipPortRepository.mu.Unlock()
	v, ok := ipPortRepository.Map[key]
	return v, ok
}

func InsertIDAndConn(key string, value *IDAndConn) {
	ipPortRepository.mu.Lock()
	defer ipPortRepository.mu.Unlock()
	ipPortRepository.Map[key] = value
}
