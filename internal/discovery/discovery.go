package discovery

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"p2p-chat/internal/types"
)

type Discovery struct {
	port            int
	username        string
	peers           map[string]string
	stopCh          chan struct{}
	mu              sync.RWMutex
	broadcastPeriod time.Duration
}

func New(port int, username string) *Discovery {
	return &Discovery{
		port:            port,
		username:        username,
		peers:           make(map[string]string),
		stopCh:          make(chan struct{}),
		broadcastPeriod: 5 * time.Second,
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
	ticker := time.NewTicker(d.broadcastPeriod)
	defer ticker.Stop()

	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+strconv.Itoa(d.port))
	if err != nil {
		fmt.Printf("Failed to resolve broadcast address: %v\n", err)
		return
	}

	for {
		select {
		case <-ticker.C:
			msg := types.NewDiscoveryMessage(d.username)
			data, err := msg.ToJSON()
			if err != nil {
				fmt.Printf("Failed to marshal discovery message: %v\n", err)
				continue
			}
			
			if _, err := conn.WriteToUDP(data, broadcastAddr); err != nil {
				fmt.Printf("Failed to send broadcast: %v\n", err)
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
			msg, err := types.MessageFromJSON(buffer[:n])
			if err != nil {
				continue
			}
			if msg.Type == types.MessageTypeDiscovery && msg.From != d.username {
				d.mu.Lock()
				d.peers[msg.From] = addr.IP.String()
				d.mu.Unlock()
				fmt.Printf("Discovered peer: %s at %s\n", msg.From, addr.IP.String())
			}
		}
	}
}

func (d *Discovery) Stop() {
	close(d.stopCh)
}

func (d *Discovery) GetPeers() map[string]string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	result := make(map[string]string)
	for name, addr := range d.peers {
		result[name] = addr
	}
	return result
}