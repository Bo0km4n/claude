package lr

import (
	"context"
	"fmt"

	"github.com/Bo0km4n/claude/app/tablet/pkg/model"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/jinzhu/gorm"
)

type LRRepository interface {
	StoreLR(ctx context.Context, lr *proto.LREntry) (*proto.LREntry, error)
	LoadLRs(ctx context.Context, lr *proto.LREntry) (*proto.LREntries, error)
	FetchLRsByDistance(ctx context.Context, latitude, longitude, distance float32) (*proto.LREntries, error)
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

func (lrr *lrRepository) FetchLRsByDistance(ctx context.Context, latitude, longitude, distance float32) (*proto.LREntries, error) {
	subQuery := fmt.Sprintf(`SELECT a.*, 
		(6371 * ACOS(COS(RADIANS(%f)) * COS(RADIANS(a.latitude)) 
		* COS(RADIANS(a.longitude) - RADIANS(%f)) + SIN(RADIANS(%f)) 
		* SIN(RADIANS(a.latitude)))) as distance from lr_entry a`, latitude, longitude, latitude)
	query := fmt.Sprintf(`SELECT id, global_ip, global_port, latitude, longitude, created_at, updated_at FROM (%s) AS d_rows WHERE d_rows.distance <= ?`, subQuery)
	rows, err := lrr.db.Raw(query, distance).Rows()
	if err != nil {
		return &proto.LREntries{}, err
	}
	defer rows.Close()

	result := &proto.LREntries{}
	for rows.Next() {
		v := &model.LREntry{}
		if err := rows.Scan(&v.ID, &v.GlobalIp, &v.GlobalPort, &v.Latitude, &v.Longitude, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return result, err
		}
		result.Entries = append(result.Entries, v.SerializeToProto())
	}

	return result, nil
}
