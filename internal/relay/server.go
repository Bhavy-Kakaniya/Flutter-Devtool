package relay

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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

func parseHandshake(message string) (command string, code string, role ClientRole) {
	parts := strings.Fields(message)
	if len(parts) == 0 {
		return "", "", ""
	}
	if len(parts) == 1 {
		role := ClientRole(parts[0])

		if role == Laptop || role == Phone {
			return "NEW", "", role
		}
		return "", "", ""
	}
	if len(parts) == 2 && parts[0] == "NEW" {
		role := ClientRole(parts[1])

		if role == Laptop || role == Phone {
			return "NEW", "", role
		}

		return "", "", ""
	}
	if len(parts) == 3 && parts[0] == "RECONNECT" {
		role := ClientRole(parts[2])

		if role == Laptop || role == Phone {
			return "RECONNECT", parts[1], role
		}

		return "", "", ""
	}
	return "", "", ""
}

// handle connection handles one client connection

func handleConnection(manager *SessionManager, connection net.Conn) {
	reader := bufio.NewReader(connection)

	// Read the client's handshake.
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read handshake:", err)
		connection.Close()
		return
	}

	// Parse the handshake into its components.
	command, code, role := parseHandshake(strings.TrimSpace(message))

	// Reject invalid handshake formats.
	if command == "" {
		fmt.Println("Invalid handshake:", message)
		connection.Close()
		return
	}

	// Handle initial pairing or reconnection.
	var session *Session

	if command == "NEW" {
		session = manager.AddClient(connection, role)
	}

	if command == "RECONNECT" {
		session = manager.ReconnectClient(connection, code, role)
	}

	// Reject clients that could not join a session.
	if session == nil {
		fmt.Println("Failed to join session:", code)
		connection.Write([]byte("ERROR\n"))
		connection.Close()
		return
	}

	fmt.Println("Client joined session:", session.ID)

	// Send the session code to the client.
	_, err = connection.Write([]byte(session.Code + "\n"))
	if err != nil {
		fmt.Println("Failed to send session code:", err)
		session.RemoveClient(connection, session.generation)
		connection.Close()
		return
	}

	// Start relaying when both clients are connected.
	if session.IsReady() {
		fmt.Println("Both clients are connected to session", session.ID)
		session.StartRelay()
	}
}
