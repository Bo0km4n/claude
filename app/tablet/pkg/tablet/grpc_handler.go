package tablet

import (
	"context"
	"errors"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/tablet/pkg/lr"
	"github.com/Bo0km4n/claude/app/tablet/pkg/util"
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

func (ts *TabletService) sendNotification(entry *proto.LREntry) {
	// TODO: implement the function that looks up some LR entries near argument's location.
}