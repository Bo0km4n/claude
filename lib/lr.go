package lib

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/peer/service"
	"google.golang.org/grpc"
)

func ConnectToLR(protocol string) {
	done := make(chan int)
	go service.LaunchGRPCService(done, protocol)
	<-done

	time.Sleep(2)
	service.UDPBcast()
	for {
		if service.IsCompletedJoinToLR {
			return
		}
		time.Sleep(1)
	}
}

func CryptedID(seed string) []byte {
	seed256 := sha256.Sum256([]byte(seed))
	lrID := make([]byte, 4)
	binary.BigEndian.PutUint32(lrID, service.RemoteLR.ID)
	dest := append(lrID, seed256[:]...)
	return dest
}

func DeserializeID(id string) []byte {
	b, _ := base64.StdEncoding.DecodeString(id)
	return b
}

func invokeLRLookUp(latitude, longitude, distance float32) ([]*proto.PeerEntry, error) {
	conn, err := grpc.Dial(service.RemoteLR.Addr + ":" + service.RemoteLR.GrpcPort)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := proto.NewLRClient(conn)

	req := &proto.LookUpPeerRequest{
		Latitude:  latitude,
		Longitude: longitude,
		Distance:  distance,
	}
	resp, err := client.LookUpPeerRPC(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.Entries, nil
}
