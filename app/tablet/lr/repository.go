package lr

import (
	"context"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/jinzhu/gorm"
)

type LRRepository interface {
	StoreLR(ctx context.Context, lr *proto.LREntry) error
	LoadLRs(ctx context.Context, lr *proto.LREntry) (*proto.LREntries, error)
}

type lrRepository struct {
	db *gorm.DB
}

func NewLRRepository(db *gorm.DB) LRRepository {
	return &lrRepository{
		db: db,
	}
}

func (lrr *lrRepository) StoreLR(ctx context.Context, in *proto.LREntry) error {
	return lrr.db.Create(in).Error
}

func (lrr *lrRepository) LoadLRs(ctx context.Context, in *proto.LREntry) (*proto.LREntries, error) {
	result := &proto.LREntries{}
	if err := lrr.db.Where(in).Find(&result.Entries).Error; err != nil {
		return nil, err
	}
	return result, nil
}
