package service

import (
	"context"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/lr/db"
)

type LRService struct{}

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
	db.InsertEntry(entry.PeerId, entry)
	// db.Dump()
	return &proto.PeerJoinResponse{Success: true}, nil
}
