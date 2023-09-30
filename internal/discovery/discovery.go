package discovery

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"p2p-chat/internal/chat"
)

type Discovery struct {
	port     int
	username string
	peers    map[string]string
	stopCh   chan struct{}
}

func New(port int, username string) *Discovery {
	return &Discovery{
		port:     port,
		username: username,
		peers:    make(map[string]string),
		stopCh:   make(chan struct{}),
	}
}

func (d *Discovery) Start() error {
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(d.port))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	go d.broadcast(conn)
	go d.listen(conn)

	<-d.stopCh
	return nil
}

func (d *Discovery) broadcast(conn *net.UDPConn) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	broadcastAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255:"+strconv.Itoa(d.port))

	for {
		select {
		case <-ticker.C:
			msg := chat.NewDiscoveryMessage(d.username)
			data, err := msg.ToJSON()
			if err == nil {
				conn.WriteToUDP(data, broadcastAddr)
			}
		case <-d.stopCh:
			return
		}
	}
}

func (d *Discovery) listen(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		select {
		case <-d.stopCh:
			return
		default:
			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				continue
			}
			msg, err := chat.MessageFromJSON(buffer[:n])
			if err != nil {
				continue
			}
			if msg.Type == chat.MessageTypeDiscovery && msg.From != d.username {
				d.peers[msg.From] = addr.IP.String()
				fmt.Printf("Discovered peer: %s at %s\n", msg.From, addr.IP.String())
			}
		}
	}
}

func (d *Discovery) Stop() {
	close(d.stopCh)
}

func (d *Discovery) GetPeers() map[string]string {
	return d.peers
}