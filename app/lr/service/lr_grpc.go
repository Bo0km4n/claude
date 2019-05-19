package service

import (
	"context"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/lr/repository"
	"google.golang.org/grpc"
)

var LRSvc *LRService

type LRService struct {
	ID uint32
}

func (p *LRService) Heartbeat(ctx context.Context, in *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, nil
}

func (p *LRService) PeerJoinRPC(ctx context.Context, in *proto.PeerJoinRequest) (*proto.PeerJoinResponse, error) {
	entry := &proto.PeerEntry{
		PeerId:    in.PeerId,
		LocalIp:   in.LocalIp,
		LocalPort: in.LocalPort,
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
	repository.InsertPeerEntry(entry.PeerId, entry)
	return &proto.PeerJoinResponse{Success: true}, nil
}

func (p *LRService) ExchangeEntriesStubRPC(ctx context.Context, in *proto.ExchangeEntriesNotification) (*proto.Empty, error) {
	localPeers := repository.FetchLocalPeers()
	req := &proto.ExchangeEntriesRequest{
		Entries: localPeers,
	}
	for _, dst := range in.Destinations {
		conn, err := grpc.Dial(dst.GlobalIp+":"+dst.GlobalPort, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := proto.NewLRClient(conn)
		resp, err := client.ExchangeEntriesDriverRPC(ctx, req)
		if err != nil {
			continue
		}
		p.registerRemotePeers(resp.Entries)
	}
	return &proto.Empty{}, nil
}

func (p *LRService) ExchangeEntriesDriverRPC(ctx context.Context, in *proto.ExchangeEntriesRequest) (*proto.ExchangeEntriesResponse, error) {
	p.registerRemotePeers(in.Entries)
	localPeers := repository.FetchLocalPeers()
	res := &proto.ExchangeEntriesResponse{
		Entries: localPeers,
	}
	return res, nil
}

func (p *LRService) registerRemotePeers(peers []*proto.PeerEntry) {
	for _, peer := range peers {
		peer.IsRemote = true
		repository.InsertPeerEntry(peer.PeerId, peer)
	}
}
