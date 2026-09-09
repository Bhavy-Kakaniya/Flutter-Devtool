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

	defer listener.Close() // close when function exists
	fmt.Println("Relay server listening on port 9000")

	var clients []net.Conn // temporary stores connected clients

	// wait for new clients
	for {
		connection, err := listener.Accept() // accept waits until client connects

		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		fmt.Println("Client connected:", connection.RemoteAddr())

		clients = append(clients, connection) // add this client to list
		fmt.Println("Connected clients:", len(clients))

		if len(clients) == 2 {
			fmt.Println("Two client connected, starting relay...")
			clientA := clients[0]
			clientB := clients[1]

			go forward(clientA, clientB)
			go forward(clientB, clientA)

			clients = nil // reset client list so another pair can be created later
		}
	}
}

// forward continuosly copies data from one connection to another
func forward(source net.Conn, destination net.Conn) {
	buffer := make([]byte, 4096) // store incoming bytes

	for {
		numberOfBytes, err := source.Read(buffer)
		if err != nil {
			fmt.Println("Connection closed:", source.RemoteAddr())
			return
		}

		_, err = destination.Write(buffer[:numberOfBytes])
		if err != nil {
			fmt.Println("Failed to forward data:", err)
			return
		}

	}
}