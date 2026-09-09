package main

import (
	"fmt"
	"net"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Failed to connect to relay:", err)
		return
	}
	defer connection.Close()

	fmt.Println("Connected to relay server")

	// writing msg to TCP connection
	_, err = connection.Write([]byte("Hello from flutter devtool"))

	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}
	fmt.Println("Message sent")
}