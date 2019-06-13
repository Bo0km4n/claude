# How to create bidirection connection

## Prepare iptables

```
sudo iptables -t nat -A POSTROUTING -s 192.168.10.100 -j SNAT --to 100.100.100.100
sudo iptables -t nat -A PREROUTING -d 100.100.100.100 -j DNAT --to 192.168.10.100
sudo iptables -t nat -A POSTROUTING -s 172.168.10.100 -j SNAT --to 100.100.100.200
sudo iptables -t nat -A PREROUTING -d 100.100.100.200 -j DNAT --to 172.168.10.100
sudo iptables -t nat -A POSTROUTING -s 200.200.200.200 -j SNAT --to 10.10.10.10 
sudo iptables -t nat -A PREROUTING -d 10.10.10.10 -j DNAT --to 200.200.200.200
```

**Peer A**: 192.168.10.101
**Proxy A**: 192.168.10.100 => 100.100.100.100(NAT Traversed)

**Peer B**: 172.168.10.100
**Proxy B**: 172.168.10.100 => 100.100.100.200(NAT Traversed)

**Router**: 192.168.10.10, 172.168.10.10, 200.200.200.10

**Tablet**: 200.200.200.200 => 10.10.10.10(NAT Traversed)

## Run example program with proxy

1. start tablet db `cd pkg/tablet && docker-compose up -d`
2. start tablet server `go run pkg/tablet/main.go`
3. start each proxy server `go run pkg/proxy/main.go --tablet_ip=10.10.10.10 --tablet_port=50051`
4. start each peer `go run examples/client/client.go udp(or tcp) xxxxxx(connect peer id)`