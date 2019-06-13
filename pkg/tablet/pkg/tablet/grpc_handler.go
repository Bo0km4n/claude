package tablet

import (
	"context"
	"errors"
	"log"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"github.com/Bo0km4n/claude/pkg/tablet/pkg/proxy"
	"github.com/Bo0km4n/claude/pkg/tablet/pkg/util"
	"google.golang.org/grpc"
)

type TabletService struct {
	proxyRepository proxy.ProxyRepository
}

func NewTabletService(proxyRepo proxy.ProxyRepository) *TabletService {
	return &TabletService{
		proxyRepository: proxyRepo,
	}
}

func (ts *TabletService) ProxyJoinRPC(ctx context.Context, in *proto.ProxyJoinRequest) (*proto.ProxyEntry, error) {
	entry := &proto.ProxyEntry{
		UniqueKey:  in.UniqueKey,
		GlobalIp:   util.GetRemoteIp(ctx),
		GlobalPort: in.GlobalPort,
		Longitude:  in.Longitude,
		Latitude:   in.Latitude,
	}
	row, err := ts.proxyRepository.StoreProxy(ctx, entry)
	if err != nil {
		return &proto.ProxyEntry{}, err
	}

	go ts.sendNotification(row)

	return row, nil
}

func (ts *TabletService) LookUpRPC(ctx context.Context, in *proto.LookUpRequest) (*proto.ProxyEntry, error) {
	query := &proto.ProxyEntry{
		Id: in.Id,
	}
	result, err := ts.proxyRepository.LoadProxys(ctx, query)
	if err != nil {
		return &proto.ProxyEntry{}, err
	}
	if len(result.Entries) <= 0 {
		return &proto.ProxyEntry{}, errors.New("Not found Proxy")
	}
	return result.Entries[0], nil
}

func (ts *TabletService) LookUpPeersRPC(ctx context.Context, in *proto.LookUpPeerRequest) (*proto.LookUpPeerResponse, error) {
	proxys, err := ts.proxyRepository.FetchProxysByDistance(ctx, in.Latitude, in.Longitude, in.Distance)
	if err != nil {
		return nil, err
	}
	peers := ts.fetchPeers(ctx, proxys.Entries)
	return &proto.LookUpPeerResponse{
		Entries: peers,
	}, nil
}

func (ts *TabletService) fetchPeers(ctx context.Context, proxys []*proto.ProxyEntry) []*proto.PeerEntry {
	peers := []*proto.PeerEntry{}
	for _, proxy := range proxys {
		conn, err := grpc.Dial(proxy.GlobalIp+":"+proxy.GlobalPort, grpc.WithInsecure())
		if err != nil {
			continue
		}
		defer conn.Close()
		client := proto.NewProxyClient(conn)
		resp, err := client.FetchPeersRPC(ctx, &proto.FetchPeersRequest{})
		if err != nil {
			log.Printf("Failed fetch peers from %s: %v", proxy.GlobalIp+":"+proxy.GlobalPort, err)
			continue
		}
		peers = append(peers, resp.Entries...)
	}
	return peers
}

// sendNotification sends notification Proxy nodes neer by argument's entry.
func (ts *TabletService) sendNotification(entry *proto.ProxyEntry) {
	distance := float32(5.0) // FIXME: This distance setting is temporary. We should modify to be able to operational.
	candidates, err := ts.proxyRepository.FetchProxysByDistance(context.Background(), entry.Latitude, entry.Longitude, distance)
	if err != nil {
		log.Println(err)
		return
	}
	if err := ts.stubExchangeEntries(context.Background(), entry, candidates); err != nil {
		log.Println(err)
		return
	}
}

func (ts *TabletService) stubExchangeEntries(ctx context.Context, newbie *proto.ProxyEntry, candidates *proto.ProxyEntries) error {
	conn, err := grpc.Dial(newbie.GlobalIp+":"+newbie.GlobalPort, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := proto.NewProxyClient(conn)
	_, err = client.ExchangeEntriesStubRPC(ctx, &proto.ExchangeEntriesNotification{
		Destinations: candidates.Entries,
	})
	return err
}
