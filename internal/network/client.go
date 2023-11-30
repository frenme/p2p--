package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"p2p-chat/internal/types"
)

const (
	DefaultMaxRetries  = 3
	DefaultRetryDelay  = time.Second
	DefaultConnTimeout = 5 * time.Second
)

type TCPClient struct {
	maxRetries  int
	retryDelay  time.Duration
	connTimeout time.Duration
}

func NewTCPClient() *TCPClient {
	return &TCPClient{
		maxRetries:  DefaultMaxRetries,
		retryDelay:  DefaultRetryDelay,
		connTimeout: DefaultConnTimeout,
	}
}

func (c *TCPClient) SendMessage(address string, port int, msg *types.Message) error {
	var lastErr error

	for i := 0; i < c.maxRetries; i++ {
		if i > 0 {
			time.Sleep(c.retryDelay)
		}

		lastErr = c.sendMessageOnce(address, port, msg)
		if lastErr == nil {
			return nil
		}
	}

	return fmt.Errorf("failed to send message after %d retries: %v", c.maxRetries, lastErr)
}

func (c *TCPClient) sendMessageOnce(address string, port int, msg *types.Message) error {
	addr := fmt.Sprintf("%s:%d", address, port)
	conn, err := net.DialTimeout("tcp", addr, c.connTimeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", addr, err)
	}
	defer conn.Close()

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	data = append(data, '\n')
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to write message to %s: %v", addr, err)
	}

	return nil
}
