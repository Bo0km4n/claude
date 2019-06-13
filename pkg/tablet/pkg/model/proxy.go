package model

import (
	"time"

	"github.com/Bo0km4n/claude/pkg/common/proto"
)

type ProxyEntry struct {
	ID         uint32    `gorm:"primary_key"`
	UniqueKey  string    `json:"unique_key" gorm:"unique_index;not_null"`
	GlobalIp   string    `json:"global_ip,omitempty"`
	GlobalPort string    `json:"global_port,omitempty"`
	Latitude   float32   `json:"latitude,omitempty"`
	Longitude  float32   `json:"longtitude,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

func (proxy *ProxyEntry) ParseProto(in *proto.ProxyEntry) {
	proxy.ID = in.Id
	proxy.GlobalIp = in.GlobalIp
	proxy.GlobalPort = in.GlobalPort
	proxy.Latitude = in.Latitude
	proxy.Longitude = in.Longitude
	proxy.UniqueKey = in.UniqueKey

	if in.CreatedAt != 0 {
		proxy.CreatedAt = time.Unix(in.CreatedAt, 0)
	}
	if in.UpdatedAt != 0 {
		proxy.UpdatedAt = time.Unix(in.UpdatedAt, 0)
	}
}

func (proxy *ProxyEntry) SerializeToProto() *proto.ProxyEntry {
	return &proto.ProxyEntry{
		Id:         proxy.ID,
		GlobalIp:   proxy.GlobalIp,
		GlobalPort: proxy.GlobalPort,
		Longitude:  proxy.Longitude,
		Latitude:   proxy.Latitude,
		UniqueKey:  proxy.UniqueKey,
		CreatedAt:  proxy.CreatedAt.Unix(),
		UpdatedAt:  proxy.UpdatedAt.Unix(),
	}
}
