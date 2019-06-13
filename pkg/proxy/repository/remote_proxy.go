package repository

import "sync"

type RemoteProxyRepo struct {
	mu   sync.Mutex
	dict map[string]string // key: Remote Peer ID, value: remote proxy ip
}

var remoteProxyRepo *RemoteProxyRepo

func FetchRemoteProxyIP(key string) (string, bool) {
	remoteProxyRepo.mu.Lock()
	defer remoteProxyRepo.mu.Unlock()
	v, ok := remoteProxyRepo.dict[key]
	return v, ok
}
