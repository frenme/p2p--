package errors

import (
	"fmt"
)

type ChatError struct {
	Code    string
	Message string
	Cause   error
}

func (e *ChatError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *ChatError) Unwrap() error {
	return e.Cause
}

func New(code, message string) *ChatError {
	return &ChatError{
		Code:    code,
		Message: message,
	}
}

func Wrap(code, message string, cause error) *ChatError {
	return &ChatError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

var (
	ErrNetworkUnavailable = New("NET001", "Network unavailable")
	ErrPeerNotFound       = New("PEER001", "Peer not found")
	ErrInvalidMessage     = New("MSG001", "Invalid message format")
	ErrDiscoveryFailed    = New("DISC001", "Discovery service failed")
	ErrConnectionFailed   = New("CONN001", "Connection failed")
)