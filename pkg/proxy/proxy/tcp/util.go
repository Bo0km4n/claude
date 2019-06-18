package tcp

import (
	"errors"
	"net"

	"github.com/Bo0km4n/claude/pkg/proxy/repository/remotepeer"
)

func newConnectionToProxy(id string) (net.Conn, error) {
	ip, ok := remotepeer.FetchRemoteProxyIP(id)
	if !ok {
		// TODO: fetch proxy information from tablet
		return nil, errors.New("Not found proxy")
	}
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
