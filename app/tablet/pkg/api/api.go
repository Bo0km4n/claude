package api

import (
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/tablet/config"
	"github.com/Bo0km4n/claude/app/tablet/pkg/db"
	"github.com/Bo0km4n/claude/app/tablet/pkg/lr"

	"github.com/Bo0km4n/claude/app/common/proto"

	"github.com/Bo0km4n/claude/app/tablet/pkg/tablet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func GRPC() {
	port, err := net.Listen("tcp", ":"+config.Config.GRPC.Port)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()

	{
		// Tablet
		{
			lrRepo := lr.NewLRRepository(db.Mysql)
			svc := tablet.NewTabletService(lrRepo)
			proto.RegisterTabletServer(
				server,
				svc,
			)
		}
	}

	log.Println("Start grpc services...")
	reflection.Register(server)
	server.Serve(port)
}
