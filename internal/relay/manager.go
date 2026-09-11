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
	sessions map[int]*Session

	nextSessionID int // generate unique id for every newly created session

	mu sync.Mutex // mu protects session map and nextSessionId, and prevent race condition
}

// create empty session manager
func NewSessionManager() *SessionManager {

	// create manager
	return &SessionManager{
		sessions:      make(map[int]*Session),
		nextSessionID: 1,
	}
}

// AddClient adds client to an available session
// if existing session has only one client new client joins that session
// if no available session exists new session is created
// it returns the session the client joined

func (manager *SessionManager) AddClient(connection net.Conn) *Session {
	manager.mu.Lock() // lock manager as it is going to be read/write shared session state

	defer manager.mu.Unlock()

	for _, session := range manager.sessions {
		// from all existing sessions try to add client to this session
		// accessing is safe as addclient has its own mutex
		if session.AddClient(connection) {
			return session
		}
	}
	// no existing session had empty slot

	session := NewSession(manager.nextSessionID) // new session with next available id
	// session1 gets callback that this is function which should be calleed when this session need to be remove
	session.removeCallback = manager.RemoveSession // notify manager that session has been expired
	manager.nextSessionID++                        // get another id for new session
	manager.sessions[session.ID] = session         // add this session to map
	session.AddClient(connection)                  // add client to new session
	fmt.Println("Created new session:", session.ID)
	return session
}

// RemoveSession removes session from manager

func (manager *SessionManager) RemoveSession(sessionID int) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	delete(manager.sessions, sessionID) // remove session from map
	fmt.Println("Removed session:", sessionID)
}
