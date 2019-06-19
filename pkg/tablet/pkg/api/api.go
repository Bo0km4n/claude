package api

import (
	"log"
	"net"

	"github.com/Bo0km4n/claude/pkg/tablet/config"
	"github.com/Bo0km4n/claude/pkg/tablet/pkg/db"
	"github.com/Bo0km4n/claude/pkg/tablet/pkg/proxy"

	"github.com/Bo0km4n/claude/pkg/common/proto"

	"github.com/Bo0km4n/claude/pkg/tablet/pkg/tablet"
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
			proxyRepo := proxy.NewProxyRepository(db.Mysql)
			svc := tablet.NewTabletService(proxyRepo)
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
