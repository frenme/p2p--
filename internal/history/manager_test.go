package history

import (
	"os"
	"testing"

	"p2p-chat/internal/types"
)

func TestNewManager(t *testing.T) {
	manager := NewManager("testuser")
	
	if manager == nil {
		t.Fatal("Manager should not be nil")
	}
	
	if manager.filename != "chat_history_testuser.json" {
		t.Errorf("Expected filename 'chat_history_testuser.json', got '%s'", manager.filename)
	}
	
	if manager.messages == nil {
		t.Error("Messages slice should be initialized")
	}
}

func TestAddMessage(t *testing.T) {
	manager := NewManager("testuser")
	msg := types.NewTextMessage("user1", "hello")
	
	manager.AddMessage(msg)
	
	messages := manager.GetMessages()
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}
	
	if messages[0].Content != "hello" {
		t.Errorf("Expected content 'hello', got '%s'", messages[0].Content)
	}
}

func TestClear(t *testing.T) {
	manager := NewManager("testuser")
	
	manager.AddMessage(types.NewTextMessage("user1", "hello"))
	manager.AddMessage(types.NewTextMessage("user2", "world"))
	
	if len(manager.GetMessages()) != 2 {
		t.Error("Expected 2 messages before clear")
	}
	
	manager.Clear()
	
	if len(manager.GetMessages()) != 0 {
		t.Error("Expected 0 messages after clear")
	}
}

func TestSaveLoadHistory(t *testing.T) {
	testFile := "test_history.json"
	defer os.Remove(testFile)
	
	manager := NewManager("testuser")
	manager.filename = testFile
	
	msg1 := types.NewTextMessage("user1", "hello")
	msg2 := types.NewTextMessage("user2", "world")
	
	manager.AddMessage(msg1)
	manager.AddMessage(msg2)
	
	err := manager.SaveHistory()
	if err != nil {
		t.Fatalf("Failed to save history: %v", err)
	}
	
	newManager := NewManager("testuser2")
	newManager.filename = testFile
	
	err = newManager.LoadHistory()
	if err != nil {
		t.Fatalf("Failed to load history: %v", err)
	}
	
	loadedMessages := newManager.GetMessages()
	if len(loadedMessages) != 2 {
		t.Errorf("Expected 2 loaded messages, got %d", len(loadedMessages))
	}
	
	if loadedMessages[0].Content != "hello" {
		t.Errorf("Expected first message 'hello', got '%s'", loadedMessages[0].Content)
	}
}

func TestLoadNonExistentHistory(t *testing.T) {
	manager := NewManager("testuser")
	manager.filename = "non_existent_file.json"
	
	err := manager.LoadHistory()
	if err != nil {
		t.Errorf("Loading non-existent file should not error, got: %v", err)
	}
	
	if len(manager.GetMessages()) != 0 {
		t.Error("Expected empty messages for non-existent file")
	}
}