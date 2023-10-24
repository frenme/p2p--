package config

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	
	if cfg == nil {
		t.Fatal("Config should not be nil")
	}
	
	if cfg.DiscoveryPort != 8080 {
		t.Errorf("Expected DiscoveryPort 8080, got %d", cfg.DiscoveryPort)
	}
	
	if cfg.TCPPort != 8081 {
		t.Errorf("Expected TCPPort 8081, got %d", cfg.TCPPort)
	}
	
	if cfg.BroadcastPeriod != 5*time.Second {
		t.Errorf("Expected BroadcastPeriod 5s, got %v", cfg.BroadcastPeriod)
	}
	
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Errorf("Expected ShutdownTimeout 10s, got %v", cfg.ShutdownTimeout)
	}
}