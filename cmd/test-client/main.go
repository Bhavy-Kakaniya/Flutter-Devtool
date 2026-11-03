package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  go run ./cmd/test-client NEW [LAPTOP|PHONE]")
		fmt.Println("  go run ./cmd/test-client RECONNECT <session-code> [LAPTOP|PHONE]")
		return
	}

	command := os.Args[1]

	var handshake string

	if command == "NEW" {
		role := os.Args[2]
		handshake = "NEW " + role
	} else if command == "RECONNECT" {
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run ./cmd/test-client RECONNECT <session-code> [LAPTOP|PHONE]")
			return
		}

		code := os.Args[2]
		role := os.Args[3]

		handshake = "RECONNECT " + code + " " + role
	} else {
		fmt.Println("Unknown command:", command)
		return
	}

	connection, reader, sessionCode, err := connectToRelay(handshake)
	if err != nil {
		fmt.Println("Failed to connect to relay:", err)
		return
	}
	fmt.Println("Joined session:", sessionCode)

	defer connection.Close()

	fmt.Println("Connected to relay server")

	go func() {
		buffer := make([]byte, 4096)

		for {
			numberOfBytes, err := reader.Read(buffer)
			if err != nil {
				fmt.Println("Connection closed")
				return
			}

			fmt.Println("Recieved:", string(buffer[:numberOfBytes]))
			fmt.Print("> ")
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			return
		}

		message := scanner.Text()
		_, err := connection.Write([]byte(message))

		if err != nil {
			fmt.Println("Failed to send message:", err)
			return
		}
	}
}

func connectToRelay(handshake string) (net.Conn, *bufio.Reader, string, error) {
	connection, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		return nil, nil, "", err
	}

	_, err = connection.Write([]byte(handshake + "\n"))
	if err != nil {
		connection.Close()
		return nil, nil, "", err
	}
	reader := bufio.NewReader(connection)

	sessionCode, err := reader.ReadString('\n')
	if err != nil {
		connection.Close()
		return nil, nil, "", err
	}
	sessionCode = strings.TrimSpace(sessionCode)
	return connection, reader, sessionCode, nil
}
