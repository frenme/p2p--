package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DiscoveryPort   int
	TCPPort         int
	BroadcastPeriod time.Duration
	ShutdownTimeout time.Duration
}

func DefaultConfig() *Config {
	discoveryPort := 8080
	tcpPort := 8081
	
	if port := os.Getenv("P2P_DISCOVERY_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			discoveryPort = p
		}
	}
	
	if port := os.Getenv("P2P_TCP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			tcpPort = p
		}
	}
	
	return &Config{
		DiscoveryPort:   discoveryPort,
		TCPPort:         tcpPort,
		BroadcastPeriod: 5 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}
}