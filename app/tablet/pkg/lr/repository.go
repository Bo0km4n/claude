package lr

import (
	"context"

	"github.com/Bo0km4n/claude/app/tablet/pkg/model"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/jinzhu/gorm"
)

type LRRepository interface {
	StoreLR(ctx context.Context, lr *proto.LREntry) (*proto.LREntry, error)
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

func (lrr *lrRepository) StoreLR(ctx context.Context, in *proto.LREntry) (*proto.LREntry, error) {
	query := &model.LREntry{}
	query.ParseProto(in)
	if err := lrr.db.Create(query).Error; err != nil {
		return &proto.LREntry{}, err
	}
	return query.SerializeToProto(), nil
}

func (lrr *lrRepository) LoadLRs(ctx context.Context, in *proto.LREntry) (*proto.LREntries, error) {
	result := []model.LREntry{}
	query := &model.LREntry{}
	query.ParseProto(in)

	if err := lrr.db.Where(query).Find(&result).Error; err != nil {
		return nil, err
	}

	protoResult := &proto.LREntries{}
	for i := range result {
		protoResult.Entries = append(protoResult.Entries, result[i].SerializeToProto())
	}
	return protoResult, nil
}
