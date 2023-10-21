package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"p2p-chat/internal/chat"
	"p2p-chat/internal/cli"
	"p2p-chat/internal/config"
	"p2p-chat/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <username>")
		return
	}
	
	username := os.Args[1]
	cfg := config.DefaultConfig()
	
	c := chat.NewChat(username, cfg.DiscoveryPort, cfg.TCPPort)
	shutdown := utils.NewGracefulShutdown(cfg.ShutdownTimeout)
	
	go func() {
		shutdown.WaitForSignal()
		c.Stop()
	}()
	
	go func() {
		if err := c.Start(); err != nil {
			fmt.Printf("Error starting chat: %v\n", err)
		}
	}()
	
	fmt.Println("Type /help for available commands")
	
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		
		isCommand, err := cli.HandleCommand(c, input)
		if isCommand {
			if err != nil {
				if err.Error() == "quit" {
					break
				}
				fmt.Printf("Error: %v\n", err)
			}
			continue
		}
		
		c.SendMessageToPeers(input)
	}
	
	c.Stop()
}