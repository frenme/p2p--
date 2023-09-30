package network

import (
	"net"
	"strconv"
)

type TCPServer struct {
	port     int
	listener net.Listener
	stopCh   chan struct{}
}

func NewTCPServer(port int) *TCPServer {
	return &TCPServer{
		port:   port,
		stopCh: make(chan struct{}),
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
}

func (s *TCPServer) Stop() {
	close(s.stopCh)
	if s.listener != nil {
		s.listener.Close()
	}
}