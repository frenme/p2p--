package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	mu                sync.RWMutex
	messagesSent      uint64
	messagesReceived  uint64
	peersDiscovered   uint64
	connectionErrors  uint64
	startTime         time.Time
}

var global = &Metrics{
	startTime: time.Now(),
}

func Global() *Metrics {
	return global
}

func (m *Metrics) IncrementMessagesSent() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messagesSent++
}

func (m *Metrics) IncrementMessagesReceived() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messagesReceived++
}

func (m *Metrics) IncrementPeersDiscovered() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.peersDiscovered++
}

func (m *Metrics) IncrementConnectionErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connectionErrors++
}

func (m *Metrics) GetStats() (uint64, uint64, uint64, uint64, time.Duration) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	uptime := time.Since(m.startTime)
	return m.messagesSent, m.messagesReceived, m.peersDiscovered, m.connectionErrors, uptime
}

func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.messagesSent = 0
	m.messagesReceived = 0
	m.peersDiscovered = 0
	m.connectionErrors = 0
	m.startTime = time.Now()
}