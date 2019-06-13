package proxy

import (
	"context"
	"fmt"

	"github.com/Bo0km4n/claude/pkg/tablet/pkg/model"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"github.com/jinzhu/gorm"
)

type ProxyRepository interface {
	StoreProxy(ctx context.Context, proxy *proto.ProxyEntry) (*proto.ProxyEntry, error)
	LoadProxys(ctx context.Context, proxy *proto.ProxyEntry) (*proto.ProxyEntries, error)
	FetchProxysByDistance(ctx context.Context, latitude, longitude, distance float32) (*proto.ProxyEntries, error)
}

type proxyRepository struct {
	db *gorm.DB
}

func NewProxyRepository(db *gorm.DB) ProxyRepository {
	return &proxyRepository{
		db: db,
	}
}

func (proxyr *proxyRepository) StoreProxy(ctx context.Context, in *proto.ProxyEntry) (*proto.ProxyEntry, error) {
	query := &model.ProxyEntry{}

	// search existince
	query.UniqueKey = in.UniqueKey
	if err := proxyr.db.Where("unique_key = ?", query.UniqueKey).First(query).Error; err != nil {
		query.ParseProto(in)

		if err := proxyr.db.Create(query).Error; err != nil {
			return &proto.ProxyEntry{}, err
		}
		return query.SerializeToProto(), nil
	}

	// update exist row
	query.GlobalIp = in.GlobalIp
	query.GlobalPort = in.GlobalPort
	query.Latitude = in.Latitude
	query.Longitude = in.Longitude
	if err := proxyr.db.Save(query).Error; err != nil {
		return &proto.ProxyEntry{}, err
	}
	return query.SerializeToProto(), nil
}

func (proxyr *proxyRepository) LoadProxys(ctx context.Context, in *proto.ProxyEntry) (*proto.ProxyEntries, error) {
	result := []model.ProxyEntry{}
	query := &model.ProxyEntry{}
	query.ParseProto(in)

	if err := proxyr.db.Where(query).Find(&result).Error; err != nil {
		return nil, err
	}

	protoResult := &proto.ProxyEntries{}
	for i := range result {
		protoResult.Entries = append(protoResult.Entries, result[i].SerializeToProto())
	}
	return protoResult, nil
}

func (proxyr *proxyRepository) FetchProxysByDistance(ctx context.Context, latitude, longitude, distance float32) (*proto.ProxyEntries, error) {
	subQuery := fmt.Sprintf(`SELECT a.*, 
		(6371 * ACOS(COS(RADIANS(%f)) * COS(RADIANS(a.latitude)) 
		* COS(RADIANS(a.longitude) - RADIANS(%f)) + SIN(RADIANS(%f)) 
		* SIN(RADIANS(a.latitude)))) as distance from proxy_entry a`, latitude, longitude, latitude)
	query := fmt.Sprintf(`SELECT id, global_ip, global_port, latitude, longitude, created_at, updated_at FROM (%s) AS d_rows WHERE d_rows.distance <= ?`, subQuery)
	rows, err := proxyr.db.Raw(query, distance).Rows()
	if err != nil {
		return &proto.ProxyEntries{}, err
	}
	defer rows.Close()

	result := &proto.ProxyEntries{}
	for rows.Next() {
		v := &model.ProxyEntry{}
		if err := rows.Scan(&v.ID, &v.GlobalIp, &v.GlobalPort, &v.Latitude, &v.Longitude, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return result, err
		}
		result.Entries = append(result.Entries, v.SerializeToProto())
	}

	return result, nil
}
