package service

import (
	"context"
	"log"
	"time"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"github.com/Bo0km4n/claude/pkg/proxy/repository/pipe"
	"github.com/Bo0km4n/claude/pkg/proxy/repository/remotepeer"
	"google.golang.org/grpc"
)

var ProxySvc *ProxyService

// ProxyService has the same information pkg/tablet/model/proxy_entry.go
type ProxyService struct {
	ID         uint32    `gorm:"primary_key"`
	GlobalIp   string    `json:"global_ip,omitempty"`
	GlobalPort string    `json:"global_port,omitempty"`
	Latitude   float32   `json:"latitude,omitempty"`
	Longitude  float32   `json:"longtitude,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

func (p *ProxyService) Heartbeat(ctx context.Context, in *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, nil
}

func (p *ProxyService) ExchangeEntriesStubRPC(ctx context.Context, in *proto.ExchangeEntriesNotification) (*proto.Empty, error) {
	localPeers := pipe.FetchLocalPeers()
	protoPeers := []*proto.PeerEntry{}
	for _, v := range localPeers {
		peer := &proto.PeerEntry{
			PeerId:    v.ID,
			ProxyIp:   p.GlobalIp + ":" + p.GlobalPort,
			Latitude:  v.Latitude,
			Longitude: v.Longitude,
		}
		protoPeers = append(protoPeers, peer)
	}
	req := &proto.ExchangeEntriesRequest{
		Entries: protoPeers,
	}
	for _, dst := range in.Destinations {
		conn, err := grpc.Dial(dst.GlobalIp+":"+dst.GlobalPort, grpc.WithInsecure())
		if err != nil {
			log.Printf("Connection create failed: %v", err)
			return nil, err
		}
		defer conn.Close()
		client := proto.NewProxyClient(conn)
		resp, err := client.ExchangeEntriesDriverRPC(ctx, req)
		if err != nil {
			log.Printf("client.ExchangeEntriesDriverRPC: %v", err)
			continue
		}
		p.registerRemotePeers(resp.Entries)
	}
	return &proto.Empty{}, nil
}

func (p *ProxyService) ExchangeEntriesDriverRPC(ctx context.Context, in *proto.ExchangeEntriesRequest) (*proto.ExchangeEntriesResponse, error) {
	p.registerRemotePeers(in.Entries)
	localPeers := pipe.FetchLocalPeers()
	protoPeers := []*proto.PeerEntry{}
	for _, v := range localPeers {
		peer := &proto.PeerEntry{
			PeerId:    v.ID,
			ProxyIp:   p.GlobalIp + ":" + p.GlobalPort,
			Latitude:  v.Latitude,
			Longitude: v.Longitude,
		}
		protoPeers = append(protoPeers, peer)
	}

	res := &proto.ExchangeEntriesResponse{
		Entries: protoPeers,
	}
	return res, nil
}

func (p *ProxyService) registerRemotePeers(peers []*proto.PeerEntry) {
	for _, peer := range peers {
		remotepeer.InsertRemotePeer(peer.PeerId, peer.ProxyIp)
	}
}

// Fetch peer entires from in-memory hash table
func (p *ProxyService) FetchPeersRPC(ctx context.Context, in *proto.FetchPeersRequest) (*proto.FetchPeersResponse, error) {
	localPeers := pipe.FetchLocalPeers()
	protoPeers := []*proto.PeerEntry{}
	for _, v := range localPeers {
		peer := &proto.PeerEntry{
			PeerId:    v.ID,
			ProxyIp:   p.GlobalIp + ":" + p.GlobalPort,
			Latitude:  v.Latitude,
			Longitude: v.Longitude,
		}
		protoPeers = append(protoPeers, peer)
	}
	return &proto.FetchPeersResponse{
		Entries: protoPeers,
	}, nil
}

// Fetch peer entries by location query via Tablet server
func (p *ProxyService) LookUpPeersRPC(ctx context.Context, in *proto.LookUpPeerRequest) (*proto.LookUpPeerResponse, error) {
	// forwarded query to Tablet.
	client, err := newTabletClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.LookUpPeersRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
