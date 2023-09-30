package chat

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	MessageTypeText      MessageType = "text"
	MessageTypeJoin      MessageType = "join"
	MessageTypeLeave     MessageType = "leave"
	MessageTypeDiscovery MessageType = "discovery"
)

type Message struct {
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewTextMessage(from, content string) *Message {
	return &Message{
		Type:      MessageTypeText,
		From:      from,
		Content:   content,
		Timestamp: time.Now(),
	}
}

func NewDiscoveryMessage(from string) *Message {
	return &Message{
		Type:      MessageTypeDiscovery,
		From:      from,
		Timestamp: time.Now(),
	}
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func MessageFromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}