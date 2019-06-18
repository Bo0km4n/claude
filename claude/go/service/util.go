package service

import (
	"context"
	"log"
	"net"

	"github.com/Bo0km4n/claude/pkg/common/proto"
	"google.golang.org/grpc"
)

func getLocalIP(dev string) string {
	iface, err := net.InterfaceByName(dev)
	if err != nil {
		log.Fatal(err)
	}
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		log.Fatal(err)
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}
					break
				}
			}
		}
	}

	return addr.IP.String()
}

func getPeerIDString(proxy *remoteProxy, seed string) (string, error) {
	conn, err := grpc.Dial(proxy.Addr+":"+proxy.GrpcPort, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	client := proto.NewProxyClient(conn)
	resp, err := client.GeneratePeerID(context.Background(), &proto.GeneratePeerIDRequest{
		Seed: seed,
	})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}
