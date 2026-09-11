package relay

import (
	"fmt"
	"net"
)

func StartServer() error {
	listener, err := net.Listen("tcp", ":9000")
	// create a TCP listener
	// :9000 means listen on all available network interfaces use port 9000
	if err != nil {
		return fmt.Errorf("Failed to start relay server: %w", err)
	}

	defer listener.Close() // close when function exits

	manager := NewSessionManager()
	fmt.Println("Relay server listening on port 9000")

	// keep accepting clients
	for {
		connection, err := listener.Accept() // accept waits until client connects

		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		fmt.Println("Client connected:", connection.RemoteAddr())

		// handle this client in its own goroutine
		// the main server goroutine immediately goes back to Accept() and can accept other clients
		go handleConnection(manager, connection)
	}
}

// handle connection handles one client connection

func handleConnection(manager *SessionManager, connection net.Conn) {
	session := manager.AddClient(connection) // add client to appropriate session
	fmt.Println("Client joined session:", session.ID)

	if session.IsReady() {
		fmt.Println("Both clients are connected to session", session.ID)
		session.StartRelay()
	}
}