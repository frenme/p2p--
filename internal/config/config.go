package config

import (
	"time"
)

type Config struct {
	DiscoveryPort   int
	TCPPort         int
	BroadcastPeriod time.Duration
	ShutdownTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		DiscoveryPort:   8080,
		TCPPort:         8081,
		BroadcastPeriod: 5 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}
}