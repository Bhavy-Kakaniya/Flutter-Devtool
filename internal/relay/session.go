package relay

import (
	"fmt"
	"net"
	"sync"
)

// Session represents one communication channel between two clients
// These two clients will be
// Client A = developer's laptop
// Client B = Android device/agent
// currently they are only 2 tcp clients

type Session struct {
	ID      int
	clientA net.Conn
	clientB net.Conn
	mu      sync.Mutex
	// mu protects session from concurrent acess, multiple goroutines work with session at same time
	// this prevents race condition
	removeCallback func(int) // called when session is closed
	closed         bool
}

func NewSession(id int) *Session {
	return &Session{
		ID: id,
	}
}

func (s *Session) AddClient(connection net.Conn) bool {
	s.mu.Lock()         // lock session before modifying its fields
	defer s.mu.Unlock() // unlock when function finishes

	if s.clientA == nil {
		s.clientA = connection
		fmt.Println("client joined session", s.ID, "as client A")
		return true
	}
	if s.clientB == nil {
		s.clientB = connection
		fmt.Println("client joined session", s.ID, "as client B")
		return true
	}
	return false
}

func (s *Session) IsReady() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clientA != nil && s.clientB != nil
}

func (s *Session) StartRelay() {
	s.mu.Lock() // lock while reading client interface
	clientA := s.clientA
	clientB := s.clientB
	s.mu.Unlock()

	go s.forward(clientA, clientB)
	go s.forward(clientB, clientA)
	fmt.Println("Session", s.ID, "is now relaying")
}

func (s *Session) Close() {
	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true

	if s.clientA != nil {
		s.clientA.Close()
		s.clientA = nil
	}
	if s.clientB != nil {
		s.clientB.Close()
		s.clientB = nil
	}

	// save callback locally
	// do this while holding lock so another goroutine cannot change callback while is is being used
	removeCallback := s.removeCallback

	s.mu.Unlock()
	fmt.Println("Session", s.ID, "closed")

	// notify SessionManager after releasing session lock
	if removeCallback != nil {
		removeCallback(s.ID)
	}
}

func (s *Session) forward(source net.Conn, destination net.Conn) {
	buffer := make([]byte, 4096) // store incoming bytes

	for {
		numberOfBytes, err := source.Read(buffer)

		// close both source and destination
		if err != nil {
			fmt.Println("Connection closed:", source.RemoteAddr())
			s.Close()
			return
		}

		_, err = destination.Write(buffer[:numberOfBytes])
		if err != nil {
			fmt.Println("Failed to forward data:", err)
			s.Close()
			return
		}
	}
}
