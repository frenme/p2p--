package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"p2p-chat/internal/chat"
	"p2p-chat/internal/cli"
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