package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"p2p-chat/internal/chat"
	"p2p-chat/internal/cli"
	"p2p-chat/internal/config"
	"p2p-chat/internal/utils"
	"p2p-chat/internal/version"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <username>")
		fmt.Printf("P2P Chat %s (Go %s)\n", version.Short(), runtime.Version())
		return
	}

	username := os.Args[1]

	if !utils.IsValidUsername(username) {
		fmt.Printf("Error: Invalid username '%s'. Username must be 2-20 characters and contain only letters, digits, underscores, and dashes.\n", username)
		return
	}

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
