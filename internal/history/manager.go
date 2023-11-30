package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"p2p-chat/internal/types"
)

type Manager struct {
	filename string
	messages []*types.Message
	mu       sync.RWMutex
}

func NewManager(username string) *Manager {
	filename := fmt.Sprintf("chat_history_%s.json", username)
	return &Manager{
		filename: filename,
		messages: make([]*types.Message, 0),
	}
}

func (m *Manager) LoadHistory() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := os.Stat(m.filename); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(m.filename)
	if err != nil {
		return fmt.Errorf("failed to read history file: %v", err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &m.messages); err != nil {
		return fmt.Errorf("failed to parse history file: %v", err)
	}

	return nil
}

func (m *Manager) SaveHistory() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.MarshalIndent(m.messages, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}

	dir := filepath.Dir(m.filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(m.filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %v", err)
	}

	return nil
}

func (m *Manager) AddMessage(msg *types.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, msg)

	if len(m.messages) > 1000 {
		m.messages = m.messages[100:]
	}
}

func (m *Manager) GetMessages() []*types.Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*types.Message, len(m.messages))
	copy(result, m.messages)
	return result
}

func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = m.messages[:0]
}
