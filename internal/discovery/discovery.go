package discovery

import (
	"net"
	"time"
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
	addr, err := net.ResolveUDPAddr("udp", ":"+string(rune(d.port)))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	go d.broadcast()
	go d.listen(conn)

	<-d.stopCh
	return nil
}

func (d *Discovery) broadcast() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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
			_, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				continue
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