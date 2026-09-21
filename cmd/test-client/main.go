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
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./cmd/test-client [LAPTOP|PHONE]")
		return
	}
	role := os.Args[1]

	connection, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Failed to connect to relay:", err)
		return
	}

	_, err = connection.Write([]byte(role + "\n"))
	if err != nil {
		fmt.Println("Failed to send role:", err)
		return
	}

	reader := bufio.NewReader(connection)

	sessionCode, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to recieve session code:", err)
		return
	}

	sessionCode = strings.TrimSpace(sessionCode)
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

	//read input by user
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
