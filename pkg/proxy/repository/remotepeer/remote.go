package remotepeer

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

func Fetch(key string) (*IDAndConn, bool) {
	ipPortRepository.mu.Lock()
	defer ipPortRepository.mu.Unlock()
	v, ok := ipPortRepository.Map[key]
	return v, ok
}

func Insert(key string, value *IDAndConn) {
	ipPortRepository.mu.Lock()
	defer ipPortRepository.mu.Unlock()
	ipPortRepository.Map[key] = value
}

func AllEntries() []*IDAndConn {
	ipPortRepository.mu.Lock()
	defer ipPortRepository.mu.Unlock()
	r := []*IDAndConn{}
	for _, v := range ipPortRepository.Map {
		r = append(r, v)
	}
	return r
}
