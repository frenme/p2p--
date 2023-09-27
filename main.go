package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <username>")
		return
	}
	
	username := os.Args[1]
	fmt.Printf("Starting P2P chat for user: %s\n", username)
}