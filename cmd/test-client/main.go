package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Failed to connect to relay:", err)
		return
	}
	defer connection.Close()

	fmt.Println("Connected to relay server")

	go func() {
		buffer := make([]byte, 4096)

		for {
			numberOfBytes, err := connection.Read(buffer)
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
