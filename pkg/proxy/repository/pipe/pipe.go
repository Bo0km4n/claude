package pipe

import (
	"net"
	"sync"

	"github.com/Bo0km4n/claude/pkg/proxy/model"
)

// Instead local peer connection repository
type IDRepo struct {
	Map map[string]*Pipe // Key: ID, Value: Pipe connection
	mu  sync.Mutex
}

type Pipe struct {
	Addr            string // IP:Port
	PeerConnection  net.Conn
	ProxyConnection net.Conn
}

func InitRepo() {
	idRepo = &IDRepo{
		Map: map[string]*Pipe{},
		mu:  sync.Mutex{},
	}
}

var idRepo *IDRepo

func Insert(key string, value *Pipe) {
	idRepo.mu.Lock()
	defer idRepo.mu.Unlock()
	idRepo.Map[key] = value
}

func Fetch(key string) (*Pipe, bool) {
	idRepo.mu.Lock()
	defer idRepo.mu.Unlock()
	v, ok := idRepo.Map[key]
	return v, ok
}

func FetchLocalPeers() []*model.Peer {
	r := []*model.Peer{}
	for k, v := range idRepo.Map {
		r = append(r, &model.Peer{
			ID:   k,
			Addr: v.Addr,
			// Longitude,
			// Latitude,
		})
	}
	return r
}
