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
	ID     int
	Code   string
	Laptop net.Conn
	Phone  net.Conn
	mu     sync.Mutex
	// mu protects session from concurrent acess, multiple goroutines work with session at same time
	// this prevents race condition
	removeCallback func(int) // called when session is closed
	started        bool
	closed         bool
}

type ClientRole string

const (
	Laptop ClientRole = "LAPTOP"
	Phone  ClientRole = "PHONE"
)

func NewSession(id int) *Session {
	return &Session{
		ID:   id,
		Code: fmt.Sprintf("Session-%d", id),
	}
}

func (s *Session) AddClient(connection net.Conn, role ClientRole) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if role == Laptop {
		// can have only one laptop connection
		if s.Laptop != nil {
			return false
		}

		s.Laptop = connection
		fmt.Println("Client joined session", s.ID, "as LAPTOP")
		return true
	}

	if role == Phone {
		// can have only one phone connection
		if s.Phone != nil {
			return false
		}

		s.Phone = connection
		fmt.Println("Client joined session", s.ID, "as PHONE")
		return true
	}
	return false // Unknown role
}

func (s *Session) IsReady() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Laptop != nil && s.Phone != nil
}

func (s *Session) StartRelay() {
	s.mu.Lock() // lock while reading client interface
	if s.started {
		s.mu.Unlock()
		return
	}
	if s.Laptop == nil || s.Phone == nil {
		s.mu.Unlock()
		return
	}
	s.started = true

	Laptop := s.Laptop
	Phone := s.Phone
	s.mu.Unlock()

	go s.forward(Laptop, Phone)
	go s.forward(Phone, Laptop)
	fmt.Println("Session", s.ID, "is now relaying")
}

func (s *Session) Close() {
	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true

	if s.Laptop != nil {
		s.Laptop.Close()
		s.Laptop = nil
	}
	if s.Phone != nil {
		s.Phone.Close()
		s.Phone = nil
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
