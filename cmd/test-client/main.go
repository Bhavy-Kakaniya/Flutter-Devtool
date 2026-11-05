package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type connectionManager struct {
	conn   net.Conn
	reader *bufio.Reader
	mu     sync.Mutex
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  go run ./cmd/test-client NEW [LAPTOP|PHONE]")
		fmt.Println("  go run ./cmd/test-client RECONNECT <session-code> [LAPTOP|PHONE]")
		return
	}

	command := os.Args[1]
	var handshake string
	var role string
	if command == "NEW" {
		role = os.Args[2]
		handshake = "NEW " + role
	} else if command == "RECONNECT" {
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run ./cmd/test-client RECONNECT <session-code> [LAPTOP|PHONE]")
			return
		}
		code := os.Args[2]
		role = os.Args[3]
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

	manager := &connectionManager{}
	manager.Replace(connection, reader)

	defer connection.Close()
	fmt.Println("Connected to relay server")

	go func() {
		for {
			currentReader := manager.Reader()
			message, err := currentReader.ReadString('\n')
			if err != nil {
				fmt.Println("\nConnection lost. Attempting to reconnect...")
				newConnection, newReader, reconnectErr := reconnectToRelay(sessionCode, role)
				if reconnectErr != nil {
					fmt.Println("Reconnection failed:", reconnectErr)
					return
				}
				manager.Replace(newConnection, newReader)
				fmt.Println("Connection restored.")
				continue
			}
			message = strings.TrimRight(message, "\r\n")
			fmt.Println("Received:", message)
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
		if message == "/disconnect" {
			manager.Disconnect()
			fmt.Println("Simulated network disconnection")
			continue
		}
		err := manager.Send(message)
		if err != nil {
			fmt.Println("Failed to send message:", err)
			continue
		}
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
	if sessionCode == "ERROR" {
		connection.Close()
		return nil, nil, "", fmt.Errorf("server rejected the connection")
	}
	return connection, reader, sessionCode, nil
}

func reconnectToRelay(sessionCode string, role string) (net.Conn, *bufio.Reader, error) {
	for {
		handshake := "RECONNECT " + sessionCode + " " + role
		connection, reader, _, err := connectToRelay(handshake)
		if err == nil {
			fmt.Println("Reconnected to session:", sessionCode)
			return connection, reader, nil
		}
		fmt.Println("Reconnect failed:", err)
		time.Sleep(2 * time.Second)
	}
}

func (manager *connectionManager) Replace(connection net.Conn, reader *bufio.Reader) {
	manager.mu.Lock()
	oldConnection := manager.conn
	manager.conn = connection
	manager.reader = reader
	manager.mu.Unlock()
	if oldConnection != nil && oldConnection != connection {
		oldConnection.Close()
	}
}

func (manager *connectionManager) Send(message string) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if manager.conn == nil {
		return fmt.Errorf("no active connection")
	}
	message = message + "\n"
	_, err := manager.conn.Write([]byte(message))
	return err
}

func (manager *connectionManager) Reader() *bufio.Reader {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	return manager.reader
}

func (manager *connectionManager) Disconnect() {
	manager.mu.Lock()
	connection := manager.conn
	manager.conn = nil
	manager.reader = nil

	manager.mu.Unlock()
	if connection != nil {
		connection.Close()
	}
}
