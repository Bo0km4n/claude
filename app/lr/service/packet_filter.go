package service

import (
	"errors"
	"log"
	"net"

	"github.com/Bo0km4n/claude/app/lr/db"
	"github.com/Bo0km4n/claude/lib"

	"github.com/Bo0km4n/claude/app/lr/config"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/k0kubun/pp"
)

func LaunchPacketFilter() {

	// TODO: Implement the filter by using the GoPacket.

	// Debug tcp listener
	go func() {
		listen, err := net.Listen("tcp", ":"+config.Config.Claude.TcpPort)
		if err != nil {
			log.Fatal(err)
		}
		defer listen.Close()

		buf := make([]byte, 1024)
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
			}
			db.RegisterConnection(conn.RemoteAddr().String(), conn)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Received packets: %d\n", n)
			// conn.Write([]byte("hello world"))
		}
	}()

	// // Debug udp listener
	// go func() {
	// 	conn, _ := net.ListenPacket("udp", ":"+config.Config.Claude.UdpPort)
	// 	defer conn.Close()

	// 	buffer := make([]byte, 1024)
	// 	for {
	// 		length, remoteAddr, _ := conn.ReadFrom(buffer)
	// 		fmt.Printf("Received from %v: %v\n", remoteAddr, buffer[:length])
	// 	}
	// }()

	go filterPacket()
}

func filterPacket() {
	iface, err := getInterface(config.Config.Interface)
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
			// switch {
			// case tcpRecvFilter(config.Config.Claude.TcpPort, packet):
			// 	pp.Println("Received tcp packet")
			// case udpRecvFilter(config.Config.Claude.UdpPort, packet):
			// 	pp.Println("Received udp packet")
			// }
			forward(packet)
		}
	}
}

var SrcIP string

func ipRecvFilter(addr *net.IPNet, packet gopacket.Packet) bool {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return false
	}
	ipv4 := ipLayer.(*layers.IPv4)
	dstIP := ipv4.DstIP.String()
	if dstIP == addr.IP.String() {
		// log.Printf("IP is src: %v, dst: %v\n", ipv4.SrcIP.String(), ipv4.DstIP.String())
		return true
	}
	SrcIP = ipv4.SrcIP.String()
	return false
}

func tcpRecvFilter(port string, packet gopacket.Packet) []byte {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	tcp := tcpLayer.(*layers.TCP)
	dstPort := tcp.DstPort.String()
	if dstPort == port && len(tcp.Payload) > 0 {
		log.Printf("TCP Port is src: %v, dst: %v\n", tcp.SrcPort.String(), dstPort)
		pp.Printf("ACK: %v, PSH: %v, SYN: %v, FIN: %v\n", tcp.ACK, tcp.PSH, tcp.SYN, tcp.FIN)
		pp.Println(len(tcp.Payload))
		return tcp.Payload
	}
	return []byte{}
}

func udpRecvFilter(port string, packet gopacket.Packet) []byte {
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	udp := udpLayer.(*layers.UDP)
	dstPort := udp.DstPort.String()
	pp.Println(udp.Payload)
	if dstPort == port && len(udp.Payload) > 0 {
		log.Printf("UDP Port is src: %v, dst: %v\n", udp.SrcPort.String(), dstPort)
		return udp.Payload
	}
	return []byte{}
}

func isTcpPacket(packet gopacket.Packet) bool {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	return tcpLayer != nil
}

func isUdpPacket(packet gopacket.Packet) bool {
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	return udpLayer != nil
}

func forward(packet gopacket.Packet) {
	var payload []byte
	if isTcpPacket(packet) {
		payload = tcpRecvFilter(config.Config.Claude.TcpPort, packet)
	} else if isUdpPacket(packet) {
		payload = udpRecvFilter(config.Config.Claude.UdpPort, packet)
	}

	if len(payload) > 0 {
		forwardPayload(payload)
	}
}

func forwardPayload(payload []byte) {
	claudePacket, err := lib.ParseHeader(payload)
	if err != nil {
		log.Println(err)
		return
	}
	pp.Println(claudePacket)
}
