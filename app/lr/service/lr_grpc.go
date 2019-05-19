package service

import (
	"context"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/lr/repository"
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
	return nil, nil
}

func (p *LRService) ExchangeEntriesDriverRPC(ctx context.Context, in *proto.ExchangeEntriesRequest) (*proto.ExchangeEntriesResponse, error) {
	return nil, nil
}
