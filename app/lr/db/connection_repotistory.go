package db

import (
	"errors"
	"net"
	"sync"
)

var connectionRepository sync.Map

func RegisterConnection(key string, value net.Conn) {
	connectionRepository.Store(key, value)
}

func LoadConnection(key string) (net.Conn, error) {
	v, ok := connectionRepository.Load(key)
	if !ok {
		return nil, errors.New("Not found connection")
	}
	return v.(net.Conn), nil
}
