package chat

import (
	"fmt"
	"sync"
	"time"

	"p2p-chat/internal/discovery"
	"p2p-chat/internal/network"
	"p2p-chat/internal/types"
)

type Chat struct {
	username  string
	discovery *discovery.Discovery
	tcpServer *network.TCPServer
	messages  []*types.Message
	mu        sync.RWMutex
	stopCh    chan struct{}
}

func NewChat(username string, discoveryPort, tcpPort int) *Chat {
	return &Chat{
		username:  username,
		discovery: discovery.New(discoveryPort, username),
		tcpServer: network.NewTCPServer(tcpPort),
		messages:  make([]*types.Message, 0),
		stopCh:    make(chan struct{}),
	}
}

func (c *Chat) Start() error {
	fmt.Printf("Starting chat for %s...\n", c.username)

	go func() {
		if err := c.discovery.Start(); err != nil {
			fmt.Printf("Discovery error: %v\n", err)
		}
	}()

	go func() {
		if err := c.tcpServer.Start(); err != nil {
			fmt.Printf("TCP server error: %v\n", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Chat started successfully!")
	fmt.Println("Type messages to send, or 'quit' to exit")

	<-c.stopCh
	return nil
}

func (c *Chat) Stop() {
	c.discovery.Stop()
	c.tcpServer.Stop()
	close(c.stopCh)
}

func (c *Chat) AddMessage(msg *types.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, msg)
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp.Format("15:04:05"), msg.From, msg.Content)
}

func (c *Chat) GetMessages() []*types.Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.messages
}

func (c *Chat) GetPeers() map[string]string {
	return c.discovery.GetPeers()
}