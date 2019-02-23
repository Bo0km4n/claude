package main

import (
	"errors"
	"flag"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	mode       = flag.String("m", "server", "command mode")
	device     = flag.String("d", "eth0", "network device")
	listenPort = flag.String("lp", "50000", "select listen port")
)

func init() {
	flag.Parse()
}

func main() {
	switch *mode {
	case "client":
		client()
	case "server":
		server()
	}
}

func client() {

}

func server() {
	iface, err := getInterface(*device)
	if err != nil {
		log.Fatal(err)
	}
	if err := scan(&iface); err != nil {
		log.Fatal(err)
	}
}

func getInterface(device string) (net.Interface, error) {
	// Get a list of all interfaces.
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		if iface.Name == device {
			return iface, nil
		}
	}
	return net.Interface{}, errors.New("Not found network interface")
}

func scan(iface *net.Interface) error {
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return err
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
	// Sanity-check that the interface has a good address.
	if addr == nil {
		return errors.New("no good IP network found")
	} else if addr.IP[0] == 127 {
		return errors.New("skipping localhost")
	} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return errors.New("mask means network is too large")
	}
	log.Printf("Using network range %v for interface %v", addr, iface.Name)

	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	// Read packet data
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		var packet gopacket.Packet
		select {
		case packet = <-in:
			if !ipRecvFilter(addr, packet) {
				continue
			}
			if !tcpRecvFilter(*listenPort, packet) {
				continue
			}
		}
	}
}

func ipRecvFilter(addr *net.IPNet, packet gopacket.Packet) bool {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		log.Println("This packet is not matched ip v4")
		return false
	}
	ipv4 := ipLayer.(*layers.IPv4)
	dstIP := ipv4.DstIP.String()
	if dstIP == addr.IP.String() {
		log.Printf("IP is src: %v, dst: %v\n", ipv4.SrcIP.String(), ipv4.DstIP.String())
		return true
	}
	return false
}

func tcpRecvFilter(port string, packet gopacket.Packet) bool {
	log.Println("filter tcp")
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		log.Println("This packet is not matched tcp")
		return false
	}
	tcp := tcpLayer.(*layers.TCP)
	dstPort := tcp.DstPort.String()
	if dstPort == port {
		log.Printf("TCP Port is src: %v, dst: %v\n", tcp.SrcPort.String(), dstPort)
		return true
	}
	log.Println(port, dstPort)
	return false
}
