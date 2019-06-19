package remotepeer

import "sync"

type RemoteProxyRepo struct {
	mu          sync.Mutex
	dict        map[string]string // key: Remote Peer ID, value: remote proxy ip
	reverseDict map[string]string // key: Remote Proxy IP, value: Remote Peer ID
}

var remoteProxyRepo *RemoteProxyRepo

func InitRepo() {
	remoteProxyRepo = &RemoteProxyRepo{
		dict: map[string]string{},
		mu:   sync.Mutex{},
	}
	ipPortRepository = &IPPortRepo{
		Map: map[string]*IDAndConn{},
		mu:  sync.Mutex{},
	}
}

func FetchRemoteProxyIP(key string) (string, bool) {
	remoteProxyRepo.mu.Lock()
	defer remoteProxyRepo.mu.Unlock()
	v, ok := remoteProxyRepo.dict[key]
	return v, ok
}

func InsertRemotePeer(key string, value string) {
	remoteProxyRepo.mu.Lock()
	defer remoteProxyRepo.mu.Unlock()
	remoteProxyRepo.dict[key] = value
}
