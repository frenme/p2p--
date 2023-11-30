// Package chat provides the main chat functionality, coordinating discovery, networking and messaging.
package chat

import (
	"fmt"
	"sync"
	"time"

	"p2p-chat/internal/discovery"
	"p2p-chat/internal/history"
	"p2p-chat/internal/metrics"
	"p2p-chat/internal/network"
	"p2p-chat/internal/types"
	"p2p-chat/internal/utils"
)

type Chat struct {
	username       string
	discovery      *discovery.Discovery
	tcpServer      *network.TCPServer
	tcpClient      *network.TCPClient
	historyManager *history.Manager
	messages       []*types.Message
	mu             sync.RWMutex
	stopCh         chan struct{}
}

func NewChat(username string, discoveryPort, tcpPort int) *Chat {
	return &Chat{
		username:       username,
		discovery:      discovery.New(discoveryPort, username),
		tcpServer:      network.NewTCPServer(tcpPort),
		tcpClient:      network.NewTCPClient(),
		historyManager: history.NewManager(username),
		messages:       make([]*types.Message, 0),
		stopCh:         make(chan struct{}),
	}
}

func (c *Chat) Start() error {
	fmt.Printf("Starting chat for %s...\n", c.username)

	if err := c.historyManager.LoadHistory(); err != nil {
		fmt.Printf("Warning: failed to load history: %v\n", err)
	}

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

	go c.handleIncomingMessages()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("✅ Chat started successfully!")
	fmt.Println("💭 Type messages to send, or 'quit' to exit")

	<-c.stopCh
	return nil
}

func (c *Chat) Stop() {
	if err := c.historyManager.SaveHistory(); err != nil {
		fmt.Printf("Warning: failed to save history: %v\n", err)
	}
	c.discovery.Stop()
	c.tcpServer.Stop()
	close(c.stopCh)
}

func (c *Chat) AddMessage(msg *types.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, msg)
	c.historyManager.AddMessage(msg)
	fmt.Printf("💬 [%s] %s: %s\n", msg.Timestamp.Format("15:04:05"), msg.From, msg.Content)
}

func (c *Chat) GetMessages() []*types.Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.messages
}

func (c *Chat) GetPeers() map[string]string {
	return c.discovery.GetPeers()
}

func (c *Chat) handleIncomingMessages() {
	for {
		select {
		case msg := <-c.tcpServer.GetMessageChannel():
			if msg.Type == types.MessageTypeText {
				metrics.Global().IncrementMessagesReceived()
				c.AddMessage(msg)
			}
		case <-c.stopCh:
			return
		}
	}
}

func (c *Chat) SendMessageToPeers(content string) {
	content = utils.SanitizeMessage(content)
	if content == "" {
		return
	}

	msg := types.NewTextMessage(c.username, content)
	peers := c.GetPeers()

	for _, addr := range peers {
		go func(address string) {
			if err := c.tcpClient.SendMessage(address, 8081, msg); err != nil {
				fmt.Printf("❌ Failed to send message to %s: %v\n", address, err)
				metrics.Global().IncrementConnectionErrors()
			} else {
				metrics.Global().IncrementMessagesSent()
			}
		}(addr)
	}

	c.AddMessage(msg)
}
