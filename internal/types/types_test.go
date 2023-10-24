package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewTextMessage(t *testing.T) {
	msg := NewTextMessage("testuser", "hello world")
	
	if msg.Type != MessageTypeText {
		t.Errorf("Expected type %v, got %v", MessageTypeText, msg.Type)
	}
	
	if msg.From != "testuser" {
		t.Errorf("Expected from 'testuser', got '%s'", msg.From)
	}
	
	if msg.Content != "hello world" {
		t.Errorf("Expected content 'hello world', got '%s'", msg.Content)
	}
	
	if msg.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestNewDiscoveryMessage(t *testing.T) {
	msg := NewDiscoveryMessage("testuser")
	
	if msg.Type != MessageTypeDiscovery {
		t.Errorf("Expected type %v, got %v", MessageTypeDiscovery, msg.Type)
	}
	
	if msg.From != "testuser" {
		t.Errorf("Expected from 'testuser', got '%s'", msg.From)
	}
	
	if msg.Content != "" {
		t.Errorf("Expected empty content, got '%s'", msg.Content)
	}
}

func TestMessageJSON(t *testing.T) {
	original := NewTextMessage("testuser", "hello")
	
	data, err := original.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	
	parsed, err := MessageFromJSON(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}
	
	if parsed.Type != original.Type {
		t.Errorf("Type mismatch: expected %v, got %v", original.Type, parsed.Type)
	}
	
	if parsed.From != original.From {
		t.Errorf("From mismatch: expected %s, got %s", original.From, parsed.From)
	}
	
	if parsed.Content != original.Content {
		t.Errorf("Content mismatch: expected %s, got %s", original.Content, parsed.Content)
	}
}

func TestMessageFromInvalidJSON(t *testing.T) {
	_, err := MessageFromJSON([]byte("invalid json"))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestAllMessageTypes(t *testing.T) {
	types := []MessageType{
		MessageTypeText,
		MessageTypeJoin,
		MessageTypeLeave,
		MessageTypeDiscovery,
	}
	
	for _, msgType := range types {
		msg := &Message{
			Type:      msgType,
			From:      "testuser",
			Content:   "test",
			Timestamp: time.Now(),
		}
		
		data, err := json.Marshal(msg)
		if err != nil {
			t.Errorf("Failed to marshal message type %v: %v", msgType, err)
		}
		
		var parsed Message
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Errorf("Failed to unmarshal message type %v: %v", msgType, err)
		}
		
		if parsed.Type != msgType {
			t.Errorf("Message type mismatch: expected %v, got %v", msgType, parsed.Type)
		}
	}
}