package pipe

import (
	"net"
	"sync"

	"github.com/Bo0km4n/claude/pkg/proxy/model"
)

// Instead local peer connection repository
type PipeRepo struct {
	Map map[string]*Pipe // Key: ID, Value: Pipe connection
	mu  sync.Mutex
}

type Pipe struct {
	Addr               string // IP:Port
	PeerConnection     net.Conn
	ProxyConnectionMap map[string]net.Conn // <destID, proxy connection>
}

func InitRepo() {
	pipeRepo = &PipeRepo{
		Map: map[string]*Pipe{},
		mu:  sync.Mutex{},
	}
}

var pipeRepo *PipeRepo

func Insert(key string, value *Pipe) {
	pipeRepo.mu.Lock()
	defer pipeRepo.mu.Unlock()
	pipeRepo.Map[key] = value
}

func Fetch(key string) (*Pipe, bool) {
	pipeRepo.mu.Lock()
	defer pipeRepo.mu.Unlock()
	v, ok := pipeRepo.Map[key]
	return v, ok
}

func FetchLocalPeers() []*model.Peer {
	r := []*model.Peer{}
	for k, v := range pipeRepo.Map {
		r = append(r, &model.Peer{
			ID:   k,
			Addr: v.Addr,
			// Longitude,
			// Latitude,
		})
	}
	return r
}
