package message

type UDPBcastMessage struct {
	ListenAddr string `json:"listen_addr"`
	ListenPort string `json:"listen_port"`
}
