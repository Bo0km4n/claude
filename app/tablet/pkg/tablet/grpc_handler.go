package tablet

import (
	"context"

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
		Longtitude: in.Longtitude,
		Latitude:   in.Latitude,
	}
	row, err := ts.lrRepository.StoreLR(ctx, entry)
	if err != nil {
		return &proto.LREntry{}, err
	}

	return row, nil
}
