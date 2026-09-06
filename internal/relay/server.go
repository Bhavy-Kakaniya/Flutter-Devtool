package relay

import (
	"fmt"
	"net"
)

func StartServer() error {
	listener, err := net.Listen("tcp", "9000")
	// create a TCP listener
	// :9000 means listen on all available network interfaces use port 9000
	if err != nil {
		return fmt.Errorf("Failed to start relay server: %w", err)
	}

	defer listener.Close() // close when function exists
	fmt.Println("Relay server listening on port 9000")

	for {
		connection, err := listener.Accept() // accept waits until client connects

		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		fmt.Println("Client connected:", connection.RemoteAddr())

		defer connection.Close()

		buffer := make([]byte, 1024) // temporary store bytes from client
		numberOfBytes, err := connection.Read(buffer)

		if err != nil {
			fmt.Println("Failed to read from client:", err)
			continue
		}
		fmt.Println("Recieved:", string(buffer[:numberOfBytes])) // convert received bytes into string
	}
}
