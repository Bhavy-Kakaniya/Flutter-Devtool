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
	fmt.Println("Relay server listening on port 9000")

	sessionId := 1 // give every session unique number

	// keep accepting clients
	for {
		connection, err := listener.Accept() // accept waits until client connects

		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		fmt.Println("Client connected:", connection.RemoteAddr())

		session := NewSession(sessionId)
		sessionId++
		session.AddClient(connection)

		fmt.Println("Waiting for second client for session", session.ID)

		secondConnection, err := listener.Accept()

		if err != nil {
			fmt.Println("Failed to accept second client", err)
			session.Close()
			continue
		}
		fmt.Println("Client connected:", secondConnection.RemoteAddr())

		session.AddClient(secondConnection)
		if session.IsReady() {
			fmt.Println("Both clients are connected to session", session.ID)
			session.StartRelay()
		}
	}
}