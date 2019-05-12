package model

import (
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
)

type LREntry struct {
	ID         uint32    `gorm:"primary_key"`
	GlobalIp   string    `json:"global_ip,omitempty"`
	GlobalPort string    `json:"global_port,omitempty"`
	Latitude   float32   `json:"latitude,omitempty"`
	Longtitude float32   `json:"longtitude,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

func (lr *LREntry) ParseProto(in *proto.LREntry) {
	lr.GlobalIp = in.GlobalIp
	lr.GlobalPort = in.GlobalPort
	lr.Latitude = in.Latitude
	lr.Longtitude = in.Longtitude

	if in.CreatedAt != 0 {
		lr.CreatedAt = time.Unix(in.CreatedAt, 0)
	}
	if in.UpdatedAt != 0 {
		lr.UpdatedAt = time.Unix(in.UpdatedAt, 0)
	}
}

func (lr *LREntry) SerializeToProto() *proto.LREntry {
	return &proto.LREntry{
		GlobalIp:   lr.GlobalIp,
		GlobalPort: lr.GlobalPort,
		Longtitude: lr.Longtitude,
		Latitude:   lr.Latitude,
		CreatedAt:  lr.CreatedAt.Unix(),
		UpdatedAt:  lr.UpdatedAt.Unix(),
	}
}
