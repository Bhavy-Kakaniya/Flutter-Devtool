package relay

import (
	"fmt"
	"net"
	"sync"
)

/*
 * what is race condition ?
 * problem in system where output depend on uncontrollable timing
 * two or more threads/processes try to access exact same data and address at same time
 * crud action requires multiple steps and they can be interupted
 */

// SessionManager tracks all active sessions
// relay server can have many sessions at same time
// Example:
//	Session 1 -> Client A + Client B
//	Session 2 -> Client A + Client B
//	Session 3 -> Client A + Client B

type SessionManager struct {
	// sessions stores every active session
	// The key is the session ID
	// Example:  1 = Session 1,  2 = Session 2
	sessions       map[int]*Session
	sessionsByCode map[string]*Session

	nextSessionID int // generate unique id for every newly created session

	mu sync.Mutex // mu protects session map and nextSessionId, and prevent race condition
}

// create empty session manager
func NewSessionManager() *SessionManager {

	// create manager
	return &SessionManager{
		sessions:       make(map[int]*Session),
		sessionsByCode: make(map[string]*Session),
		nextSessionID:  1,
	}
}

// AddClient adds client to an available session
// if existing session has only one client new client joins that session
// if no available session exists new session is created
// it returns the session the client joined
func (manager *SessionManager) AddClient(connection net.Conn, role ClientRole) *Session {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for _, session := range manager.sessions {
		// Try to add this client to an existing session
		// The role tells the session whether this is the laptop or phone
		if session.AddClient(connection, role) {
			return session
		}
	}
	// No existing session had an available slot for this role.

	session := NewSession(manager.nextSessionID)

	manager.sessionsByCode[session.Code] = session // Store the session using its stable code

	session.removeCallback = manager.RemoveSession // Give the session a callback so it can ask the manager to remove it

	manager.nextSessionID++ // Prepare the ID for the next newly created session

	manager.sessions[session.ID] = session // Store the new session by its internal ID

	session.AddClient(connection, role) // Add the client to the appropriate side of the session

	fmt.Println("Created new session:", session.ID, "Code:", session.Code)
	return session
}

// RemoveSession removes session from manager

func (manager *SessionManager) RemoveSession(sessionID int) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	delete(manager.sessions, sessionID) // remove session from map
	fmt.Println("Removed session:", sessionID)
}

func (manager *SessionManager) GetSessionByCode(code string) *Session {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	return manager.sessionsByCode[code]
}
