package repository

import "sync"

type IDRepo struct {
	Map map[string]string // Key: IP:Port, Value: ID
	mu  sync.Mutex
}
