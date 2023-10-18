package network

import (
	"testing"
	"time"

	"p2p-chat/internal/types"
)

func TestTCPServer(t *testing.T) {
	server := NewTCPServer(0)
	if server == nil {
		t.Fatal("Server should not be nil")
	}
	
	if server.messageCh == nil {
		t.Error("Message channel should be initialized")
	}
	
	server.Stop()
}

func TestTCPClient(t *testing.T) {
	client := NewTCPClient()
	if client == nil {
		t.Fatal("Client should not be nil")
	}
	
	if client.maxRetries != 3 {
		t.Errorf("Expected maxRetries to be 3, got %d", client.maxRetries)
	}
	
	if client.retryDelay != time.Second {
		t.Errorf("Expected retryDelay to be 1s, got %v", client.retryDelay)
	}
}

func TestSendMessageToNonExistentServer(t *testing.T) {
	client := NewTCPClient()
	msg := types.NewTextMessage("testuser", "hello")
	
	err := client.SendMessage("127.0.0.1", 19999, msg)
	if err == nil {
		t.Error("Expected error when sending to non-existent server")
	}
}