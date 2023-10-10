package discovery

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	d := New(8080, "testuser")
	if d == nil {
		t.Fatal("Discovery instance should not be nil")
	}
	if d.username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", d.username)
	}
	if d.port != 8080 {
		t.Errorf("Expected port 8080, got %d", d.port)
	}
	if d.peers == nil {
		t.Error("Peers map should be initialized")
	}
}

func TestGetPeers(t *testing.T) {
	d := New(8080, "testuser")
	peers := d.GetPeers()
	if len(peers) != 0 {
		t.Errorf("Expected empty peers map, got %d peers", len(peers))
	}
	
	d.peers["user1"] = "192.168.1.1"
	peers = d.GetPeers()
	if len(peers) != 1 {
		t.Errorf("Expected 1 peer, got %d peers", len(peers))
	}
	if peers["user1"] != "192.168.1.1" {
		t.Errorf("Expected peer address '192.168.1.1', got '%s'", peers["user1"])
	}
}

func TestStop(t *testing.T) {
	d := New(8080, "testuser")
	
	go func() {
		time.Sleep(10 * time.Millisecond)
		d.Stop()
	}()
	
	err := d.Start()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}