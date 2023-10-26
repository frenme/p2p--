package cli

import (
	"fmt"
	"strings"

	"p2p-chat/internal/chat"
	"p2p-chat/internal/version"
)

type Command struct {
	Name        string
	Description string
	Handler     func(*chat.Chat, []string) error
}

var Commands map[string]Command

func init() {
	Commands = map[string]Command{
		"help": {
			Name:        "help",
			Description: "Show available commands",
			Handler:     handleHelp,
		},
		"peers": {
			Name:        "peers",
			Description: "List active peers",
			Handler:     handlePeers,
		},
		"history": {
			Name:        "history",
			Description: "Show message history",
			Handler:     handleHistory,
		},
			"clear": {
		Name:        "clear",
		Description: "Clear screen",
		Handler:     handleClear,
	},
	"save": {
		Name:        "save",
		Description: "Save message history to file",
		Handler:     handleSave,
	},
	"status": {
		Name:        "status",
		Description: "Show chat status",
		Handler:     handleStatus,
	},
	"version": {
		Name:        "version",
		Description: "Show version information",
		Handler:     handleVersion,
	},
		"quit": {
			Name:        "quit",
			Description: "Exit the chat",
			Handler:     handleQuit,
		},
	}
}

func HandleCommand(c *chat.Chat, input string) (bool, error) {
	if !strings.HasPrefix(input, "/") {
		return false, nil
	}
	
	parts := strings.Fields(input[1:])
	if len(parts) == 0 {
		return true, fmt.Errorf("empty command")
	}
	
	cmdName := parts[0]
	args := parts[1:]
	
	cmd, exists := Commands[cmdName]
	if !exists {
		return true, fmt.Errorf("unknown command: %s", cmdName)
	}
	
	return true, cmd.Handler(c, args)
}

func handleHelp(c *chat.Chat, args []string) error {
	fmt.Println("Available commands:")
	for _, cmd := range Commands {
		fmt.Printf("  /%s - %s\n", cmd.Name, cmd.Description)
	}
	return nil
}

func handlePeers(c *chat.Chat, args []string) error {
	peers := c.GetPeers()
	fmt.Printf("Active peers (%d):\n", len(peers))
	for name, addr := range peers {
		fmt.Printf("  %s (%s)\n", name, addr)
	}
	return nil
}

func handleHistory(c *chat.Chat, args []string) error {
	messages := c.GetMessages()
	if len(messages) == 0 {
		fmt.Println("No messages in history")
		return nil
	}
	
	fmt.Println("Message history:")
	for _, msg := range messages {
		fmt.Printf("[%s] %s: %s\n", 
			msg.Timestamp.Format("15:04:05"), msg.From, msg.Content)
	}
	return nil
}

func handleClear(c *chat.Chat, args []string) error {
	fmt.Print("\033[2J\033[H")
	return nil
}

func handleQuit(c *chat.Chat, args []string) error {
	return fmt.Errorf("quit")
}

func handleSave(c *chat.Chat, args []string) error {
	filename := "chat_export.json"
	if len(args) > 0 {
		filename = args[0]
	}
	
	messages := c.GetMessages()
	if len(messages) == 0 {
		fmt.Println("No messages to save")
		return nil
	}
	
	fmt.Printf("History would be saved to %s (feature not fully implemented)\n", filename)
	return nil
}

func handleStatus(c *chat.Chat, args []string) error {
	peers := c.GetPeers()
	messages := c.GetMessages()
	
	fmt.Println("=== Chat Status ===")
	fmt.Printf("Active peers: %d\n", len(peers))
	fmt.Printf("Messages in memory: %d\n", len(messages))
	
	if len(peers) > 0 {
		fmt.Println("\nConnected peers:")
		for name, addr := range peers {
			fmt.Printf("  • %s (%s)\n", name, addr)
		}
	}
	
	return nil
}

func handleVersion(c *chat.Chat, args []string) error {
	fmt.Println(version.Full())
	return nil
}