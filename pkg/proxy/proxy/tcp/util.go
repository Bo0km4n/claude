package tcp

import (
	"context"
	"errors"
	"net"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"github.com/Bo0km4n/claude/pkg/proxy/repository/remotepeer"
	"google.golang.org/grpc"
)

func getPeerID(addr string) (string, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	client := proto.NewPeerClient(conn)
	resp, err := client.GetPeerID(context.Background(), &proto.Empty{})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

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
