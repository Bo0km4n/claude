package lib

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"github.com/Bo0km4n/claude/claude/go/service"
	"google.golang.org/grpc"
)

func ConnectToProxy(protocol string) {
	done := make(chan int)
	go service.LaunchGRPCService(done, protocol)
	<-done

	time.Sleep(2)
	service.UDPBcast()
	for {
		if service.IsCompletedJoinToProxy {
			return
		}
		time.Sleep(1)
	}
}

func CryptedID(seed string) []byte {
	seed256 := sha256.Sum256([]byte(seed))
	proxyID := make([]byte, 4)
	binary.BigEndian.PutUint32(proxyID, service.RemoteProxy.ID)
	dest := append(proxyID, seed256[:]...)
	return dest
}

func DeserializeID(id string) []byte {
	b, _ := base64.StdEncoding.DecodeString(id)
	return b
}

func LookUpPeers(latitude, longitude, distance float32) ([]*proto.PeerEntry, error) {
	return invokeProxyLookUp(latitude, longitude, distance)
}

func invokeProxyLookUp(latitude, longitude, distance float32) ([]*proto.PeerEntry, error) {
	conn, err := grpc.Dial(service.RemoteProxy.Addr+":"+service.RemoteProxy.GrpcPort, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := proto.NewProxyClient(conn)

	req := &proto.LookUpPeerRequest{
		Latitude:  latitude,
		Longitude: longitude,
		Distance:  distance,
	}
	resp, err := client.LookUpPeersRPC(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.Entries, nil
}
