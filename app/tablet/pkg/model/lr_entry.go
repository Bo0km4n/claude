package model

import (
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
)

type LREntry struct {
	ID         uint32    `gorm:"primary_key"`
	UniqueKey  string    `json:"unique_key" gorm:"unique_index;not_null"`
	GlobalIp   string    `json:"global_ip,omitempty"`
	GlobalPort string    `json:"global_port,omitempty"`
	Latitude   float32   `json:"latitude,omitempty"`
	Longitude  float32   `json:"longtitude,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

func (lr *LREntry) ParseProto(in *proto.LREntry) {
	lr.ID = in.Id
	lr.GlobalIp = in.GlobalIp
	lr.GlobalPort = in.GlobalPort
	lr.Latitude = in.Latitude
	lr.Longitude = in.Longitude
	lr.UniqueKey = in.UniqueKey

	if in.CreatedAt != 0 {
		lr.CreatedAt = time.Unix(in.CreatedAt, 0)
	}
	if in.UpdatedAt != 0 {
		lr.UpdatedAt = time.Unix(in.UpdatedAt, 0)
	}
}

func (lr *LREntry) SerializeToProto() *proto.LREntry {
	return &proto.LREntry{
		Id:         lr.ID,
		GlobalIp:   lr.GlobalIp,
		GlobalPort: lr.GlobalPort,
		Longitude:  lr.Longitude,
		Latitude:   lr.Latitude,
		UniqueKey:  lr.UniqueKey,
		CreatedAt:  lr.CreatedAt.Unix(),
		UpdatedAt:  lr.UpdatedAt.Unix(),
	}
}
