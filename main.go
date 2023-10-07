package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"p2p-chat/internal/chat"
	"p2p-chat/internal/types"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <username>")
		return
	}
	
	username := os.Args[1]
	
	c := chat.NewChat(username, 8080, 8081)
	
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		c.Stop()
	}()
	
	go func() {
		if err := c.Start(); err != nil {
			fmt.Printf("Error starting chat: %v\n", err)
		}
	}()
	
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" {
			break
		}
		if input == "peers" {
			peers := c.GetPeers()
			fmt.Printf("Active peers (%d):\n", len(peers))
			for name, addr := range peers {
				fmt.Printf("  %s (%s)\n", name, addr)
			}
			continue
		}
		if input != "" {
			c.SendMessageToPeers(input)
		}
	}
	
	c.Stop()
}