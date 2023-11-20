package utils

import (
	"crypto/rand"
	"fmt"
	"strings"
	"unicode"
)

func IsValidUsername(username string) bool {
	if len(username) < 2 || len(username) > 20 {
		return false
	}
	
	if !unicode.IsLetter(rune(username[0])) {
		return false
	}
	
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return false
		}
	}
	
	return true
}

func SanitizeMessage(content string) string {
	content = strings.TrimSpace(content)
	
	if len(content) > 500 {
		content = content[:500]
	}
	
	return content
}

func GenerateSessionID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	
	if maxLen <= 3 {
		return s[:maxLen]
	}
	
	return s[:maxLen-3] + "..."
}