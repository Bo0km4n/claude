#!/bin/sh
# each peer network
sudo iptables -t nat -D PREROUTING -d 100.100.100.200 -j DNAT --to-destination 172.168.10.100
sudo iptables -t nat -D POSTROUTING  -s 192.168.10.100  -j SNAT --to-source 100.100.100.100
sudo iptables -t nat -D PREROUTING -d 100.100.100.100 -j DNAT --to-destination 192.168.10.100
sudo iptables -t nat -D POSTROUTING -s 172.168.10.100  -j SNAT --to-source 100.100.100.200
# tablet network
sudo iptables -t nat -D PREROUTING -d 10.10.10.10 -j DNAT --to-destination 200.200.200.200
sudo iptables -t nat -D POSTROUTING -s 200.200.200.200 -j SNAT --to-source 10.10.10.10