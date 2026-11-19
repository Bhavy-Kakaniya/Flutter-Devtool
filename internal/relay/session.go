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
	// mu protects session from concurrent acess, multiple goroutines work with session at same time this prevents race condition
	removeCallback      func(int) // called when session is closed
	started             bool
	closed              bool
	generation          uint64
	laptopReaderStarted bool
	phoneReaderStarted  bool
	laptopMessageQueue  [][]byte
	phoneMessageQueue   [][]byte
	laptopGeneration uint64
	phoneGeneration  uint64
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
		if s.Laptop != nil {
			return false
		}
		s.Laptop = connection
		s.generation++
		s.laptopGeneration++
		s.laptopReaderStarted = false
		fmt.Println("Client joined session", s.ID, "as LAPTOP")
		return true
	}
	if role == Phone {
		if s.Phone != nil {
			return false
		}
		s.Phone = connection
		s.generation++
		s.phoneGeneration++
		s.phoneReaderStarted = false
		fmt.Println("Client joined session", s.ID, "as PHONE")
		return true
	}
	return false
}

func (s *Session) IsReady() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Laptop != nil && s.Phone != nil
}

func (s *Session) StartRelay() {
	s.mu.Lock()
	if s.Laptop == nil || s.Phone == nil {
		s.mu.Unlock()
		return
	}
	laptop := s.Laptop
	phone := s.Phone
	if !s.laptopReaderStarted {
		s.laptopReaderStarted = true
		go s.forward(laptop, Laptop, s.laptopGeneration)
	}
	if !s.phoneReaderStarted {
		s.phoneReaderStarted = true
		go s.forward(phone, Phone, s.phoneGeneration)
	}
	s.started = true
	s.mu.Unlock()
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
	removeCallback := s.removeCallback
	s.mu.Unlock()
	fmt.Println("Session", s.ID, "closed")

	if removeCallback != nil {
		removeCallback(s.ID)
	}
}

func (s *Session) forward(source net.Conn, role ClientRole, generation uint64) {
	buffer := make([]byte, 4096)
	for {
		numberOfBytes, err := source.Read(buffer)
		if err != nil {
			fmt.Println("Read failed:", err)
			s.RemoveClientByRole(source, role, generation)
			return
		}
		s.mu.Lock()
		var destination net.Conn
		if role == Laptop {
			destination = s.Phone
		} else if role == Phone {
			destination = s.Laptop
		}
		s.mu.Unlock()
		if destination == nil {
			s.mu.Lock()
			if role == Laptop {
				s.phoneMessageQueue = append(
					s.phoneMessageQueue,
					append([]byte(nil), buffer[:numberOfBytes]...),
				)
			} else if role == Phone {
				s.laptopMessageQueue = append(
					s.laptopMessageQueue,
					append([]byte(nil), buffer[:numberOfBytes]...),
				)
			}
			s.mu.Unlock()
			fmt.Println("Destination unavailable; message queued")
			continue
		}
		_, err = destination.Write(buffer[:numberOfBytes])
		if err != nil {
			fmt.Println("Failed to forward data:", err)
			destination.Close()
			continue
		}
		fmt.Printf("Forwarded %d bytes from %s\n", numberOfBytes, role)
	}
}

func (s *Session) RemoveClient(connection net.Conn, generation uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if generation != s.generation {
		return
	}

	if s.Laptop == connection {
		s.Laptop = nil
		s.laptopReaderStarted = false
		s.started = false
		fmt.Println("Laptop disconnected from session", s.ID)
	}

	if s.Phone == connection {
		s.Phone = nil
		s.phoneReaderStarted = false
		s.started = false
		fmt.Println("Phone disconnected from session", s.ID)
	}
}

func (s *Session) RemoveClientByRole(connection net.Conn, role ClientRole, generation uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if role == Laptop {
		if s.Laptop != connection || s.laptopGeneration != generation {
			return
		}
		s.Laptop = nil
		s.laptopReaderStarted = false
		s.started = false
		fmt.Println("Laptop disconnected from session", s.ID)
	}

	if role == Phone {
		if s.Phone != connection || s.phoneGeneration != generation {
			return
		}
		s.Phone = nil
		s.phoneReaderStarted = false
		s.started = false

		fmt.Println("Phone disconnected from session", s.ID)
	}
}

func (s *Session) flushQueueLocked(destination net.Conn, role ClientRole) error {
	var queue *[][]byte
	if role == Laptop {
		queue = &s.laptopMessageQueue
	} else if role == Phone {
		queue = &s.phoneMessageQueue
	} else {
		return fmt.Errorf("invalid client role")
	}
	for len(*queue) > 0 {
		message := (*queue)[0]
		for len(message) > 0 {
			numberOfBytes, err := destination.Write(message)
			message = message[numberOfBytes:]
			(*queue)[0] = message

			if len(message) == 0 {
				*queue = (*queue)[1:]
			}

			if err != nil {
				return err
			}

			if numberOfBytes == 0 {
				return fmt.Errorf("connection wrote zero bytes")
			}
		}
	}
	return nil
}

// safely flushes pending data for a reconnecting client
func (s *Session) FlushQueue(destination net.Conn, role ClientRole) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flushQueueLocked(destination, role)
}
