package utils

import (
	"strings"
	"testing"
)

func TestIsValidUsername(t *testing.T) {
	validUsernames := []string{
		"alice",
		"bob123",
		"user_name",
		"test-user",
		"a1",
		"verylongusername123",
	}
	
	for _, username := range validUsernames {
		if !IsValidUsername(username) {
			t.Errorf("Expected '%s' to be valid", username)
		}
	}
	
	invalidUsernames := []string{
		"a",                          // too short
		"verylongusernamethatexceeds", // too long
		"user@domain",                // invalid character
		"user name",                  // space
		"user.name",                  // dot
		"",                           // empty
	}
	
	for _, username := range invalidUsernames {
		if IsValidUsername(username) {
			t.Errorf("Expected '%s' to be invalid", username)
		}
	}
}

func TestSanitizeMessage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"", ""},
		{"normal message", "normal message"},
		{strings.Repeat("a", 600), strings.Repeat("a", 500)},
	}
	
	for _, test := range tests {
		result := SanitizeMessage(test.input)
		if result != test.expected {
			t.Errorf("SanitizeMessage(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestGenerateSessionID(t *testing.T) {
	id1 := GenerateSessionID()
	id2 := GenerateSessionID()
	
	if len(id1) != 16 {
		t.Errorf("Expected session ID length 16, got %d", len(id1))
	}
	
	if id1 == id2 {
		t.Error("Session IDs should be different")
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "he..."},
		{"hi", 3, "hi"},
		{"test", 4, "test"},
		{"test", 2, "te"},
		{"", 5, ""},
	}
	
	for _, test := range tests {
		result := TruncateString(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("TruncateString(%q, %d) = %q, expected %q", 
				test.input, test.maxLen, result, test.expected)
		}
	}
}