package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"github.com/Bo0km4n/claude/pkg/proxy/config"
	"github.com/Bo0km4n/claude/pkg/proxy/geo"
	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
)

var td *tabletDaemon

const (
	TABLET_SIG_SYNC_INIT = iota
	TABLET_SIG_SYNC
)

type tabletDaemon struct {
	signal       chan int
	tabletClient proto.TabletClient
}

func initDaemon() error {
	conn, err := grpc.Dial(config.Config.Tablet.IP+":"+config.Config.Tablet.Port, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := proto.NewTabletClient(conn)
	td = &tabletDaemon{
		signal:       make(chan int),
		tabletClient: client,
	}

	return nil
}

func newTabletClient() (proto.TabletClient, error) {
	conn, err := grpc.Dial(config.Config.Tablet.IP+":"+config.Config.Tablet.Port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := proto.NewTabletClient(conn)
	return client, nil
}

func (t *tabletDaemon) start() error {
	for {
		sig := <-t.signal
		switch sig {
		case TABLET_SIG_SYNC_INIT:
			err := t.syncInit()
			if err != nil {
				log.Fatal(err)
			}
		default:
			return fmt.Errorf("Not found tablet daemon's signal: %d", sig)
		}
	}
	return nil
}

func (t *tabletDaemon) syncInit() error {
	latitude, longitude := geo.GetLocation()
	resp, err := t.tabletClient.ProxyJoinRPC(context.Background(), &proto.ProxyJoinRequest{
		UniqueKey:  getMacAddr(config.Config.Interface),
		GlobalPort: config.Config.GRPC.Port,
		Latitude:   latitude,
		Longitude:  longitude,
	})
	if err != nil {
		return err
	}

	ProxySvc.ID = resp.Id
	ProxySvc.GlobalIp = resp.GlobalIp
	ProxySvc.GlobalPort = resp.GlobalPort
	ProxySvc.Latitude = resp.Latitude
	ProxySvc.Longitude = resp.Longitude

	pp.Println(ProxySvc)
	return nil
}
