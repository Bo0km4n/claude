package service

import (
	"context"
	"log"
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/lr/config"
	"github.com/Bo0km4n/claude/app/lr/repository"
	"google.golang.org/grpc"
)

var LRSvc *LRService

// LRService has the same information app/tablet/model/lr_entry.go
type LRService struct {
	ID         uint32    `gorm:"primary_key"`
	GlobalIp   string    `json:"global_ip,omitempty"`
	GlobalPort string    `json:"global_port,omitempty"`
	Latitude   float32   `json:"latitude,omitempty"`
	Longitude  float32   `json:"longtitude,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
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
		Protocol:  in.Protocol,
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
			log.Printf("Connection create failed: %v", err)
			return nil, err
		}
		defer conn.Close()
		client := proto.NewLRClient(conn)
		resp, err := client.ExchangeEntriesDriverRPC(ctx, req)
		if err != nil {
			log.Printf("client.ExchangeEntriesDriverRPC: %v", err)
			continue
		}
		p.registerRemotePeers(resp.Entries)
	}
	repository.Dump()
	return &proto.Empty{}, nil
}

func (p *LRService) ExchangeEntriesDriverRPC(ctx context.Context, in *proto.ExchangeEntriesRequest) (*proto.ExchangeEntriesResponse, error) {
	p.registerRemotePeers(in.Entries)
	localPeers := repository.FetchLocalPeers()

	duplicatedLocalPeers := make([]*proto.PeerEntry, 0)

	// rewrite peer information
	for _, peer := range localPeers {
		dupPeer := &proto.PeerEntry{
			PeerId:  peer.PeerId,
			LocalIp: LRSvc.GlobalIp,
			LocalPort: func(protocol string) string {
				switch protocol {
				case "tcp":
					return config.Config.Claude.TcpPort
				case "udp":
					return config.Config.Claude.UdpPort
				default:
					return ""
				}
			}(peer.Protocol),
			Latitude:  peer.Latitude,
			Longitude: peer.Longitude,
			Protocol:  peer.Protocol,
		}
		duplicatedLocalPeers = append(duplicatedLocalPeers, dupPeer)
	}

	res := &proto.ExchangeEntriesResponse{
		Entries: duplicatedLocalPeers,
	}
	return res, nil
}

func (p *LRService) registerRemotePeers(peers []*proto.PeerEntry) {
	for _, peer := range peers {
		peer.IsRemote = true
		repository.InsertPeerEntry(peer.PeerId, peer)
	}
}

func (p *LRService) LookUpPeerRPC(ctx context.Context, in *proto.LookUpPeerRequest) (*proto.LookUpPeerResponse, error) {
	// TODO: implement the geo location query to LR's KVS.
	// And forwarded query to Tablet.

	return nil, nil
}
