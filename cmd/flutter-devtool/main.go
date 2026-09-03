package main

import (
	"fmt"
	"os"
)

func main() {
	// check whether the user provided a command.
	if len(os.Args) < 2 {
		fmt.Println("Usage: flutter-devtool <command>")
		return
	}

	// get command entered by user
	command := os.Args[1]

	// decide which command to execute
	switch command {
	case "devices":
		runDevices()

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Available commands: devices")
	}
}