package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/lr/config"
	"github.com/Bo0km4n/claude/app/lr/geo"
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
	resp, err := t.tabletClient.LRJoinRPC(context.Background(), &proto.LRJoinRequest{
		GlobalPort: config.Config.GRPC.Port,
		Latitude:   latitude,
		Longitude:  longitude,
	})
	if err != nil {
		return err
	}

	LRSvc.ID = resp.Id
	LRSvc.GlobalIp = resp.GlobalIp
	LRSvc.GlobalPort = resp.GlobalPort
	LRSvc.Latitude = resp.Latitude
	LRSvc.Longitude = resp.Longitude

	pp.Println(LRSvc)
	return nil
}
