package types

import (
	"testing"
)

func BenchmarkNewTextMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTextMessage("testuser", "hello world")
	}
}

func BenchmarkMessageToJSON(b *testing.B) {
	msg := NewTextMessage("testuser", "hello world")
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := msg.ToJSON()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMessageFromJSON(b *testing.B) {
	msg := NewTextMessage("testuser", "hello world")
	data, _ := msg.ToJSON()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := MessageFromJSON(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}