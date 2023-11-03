// Package network provides TCP server and client functionality for P2P communication.
package network

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"

	"p2p-chat/internal/types"
)

type TCPServer struct {
	port       int
	listener   net.Listener
	stopCh     chan struct{}
	messageCh  chan *types.Message
}

func NewTCPServer(port int) *TCPServer {
	return &TCPServer{
		port:      port,
		stopCh:    make(chan struct{}),
		messageCh: make(chan *types.Message, 100),
	}
}

func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))
	if err != nil {
		return err
	}
	s.listener = listener

	go s.acceptConnections()
	<-s.stopCh
	return nil
}

func (s *TCPServer) acceptConnections() {
	for {
		select {
		case <-s.stopCh:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			go s.handleConnection(conn)
		}
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var msg types.Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			continue
		}
		
		select {
		case s.messageCh <- &msg:
		default:
		}
	}
}

func (s *TCPServer) Stop() {
	close(s.stopCh)
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *TCPServer) GetMessageChannel() <-chan *types.Message {
	return s.messageCh
}

