package tablet

import (
	"context"
	"errors"
	"log"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/tablet/pkg/lr"
	"github.com/Bo0km4n/claude/app/tablet/pkg/util"
	"google.golang.org/grpc"
)

type TabletService struct {
	lrRepository lr.LRRepository
}

func NewTabletService(lrRepo lr.LRRepository) *TabletService {
	return &TabletService{
		lrRepository: lrRepo,
	}
}

func (ts *TabletService) LRJoinRPC(ctx context.Context, in *proto.LRJoinRequest) (*proto.LREntry, error) {
	entry := &proto.LREntry{
		GlobalIp:   util.GetRemoteIp(ctx),
		GlobalPort: in.GlobalPort,
		Longitude:  in.Longitude,
		Latitude:   in.Latitude,
	}
	row, err := ts.lrRepository.StoreLR(ctx, entry)
	if err != nil {
		return &proto.LREntry{}, err
	}

	// TODO: Send notification of order that exchange information about peer entries
	// between each LR
	go ts.sendNotification(row)

	return row, nil
}

func (ts *TabletService) LookUpRPC(ctx context.Context, in *proto.LookUpRequest) (*proto.LREntry, error) {
	query := &proto.LREntry{
		Id: in.Id,
	}
	result, err := ts.lrRepository.LoadLRs(ctx, query)
	if err != nil {
		return &proto.LREntry{}, err
	}
	if len(result.Entries) <= 0 {
		return &proto.LREntry{}, errors.New("Not found LR")
	}
	return result.Entries[0], nil
}

// sendNotification sends notification LR nodes neer by argument's entry.
func (ts *TabletService) sendNotification(entry *proto.LREntry) {
	distance := float32(1.0) // FIXME: This distance setting is temporary. We should modify to be able to operational.
	candidates, err := ts.lrRepository.FetchLRsByDistance(context.Background(), entry.Latitude, entry.Longitude, distance)
	if err != nil {
		log.Println(err)
		return
	}
	if err := ts.stubExchangeEntries(context.Background(), entry, candidates); err != nil {
		log.Println(err)
		return
	}
}

func (ts *TabletService) stubExchangeEntries(ctx context.Context, newbie *proto.LREntry, candidates *proto.LREntries) error {
	conn, err := grpc.Dial(newbie.GlobalIp+":"+newbie.GlobalPort, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := proto.NewLRClient(conn)
	_, err = client.ExchangeEntriesStubRPC(ctx, &proto.ExchangeEntriesNotification{
		Destinations: candidates.Entries,
	})
	return err
}
