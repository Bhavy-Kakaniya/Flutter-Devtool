package main

import (
	"fmt"
	"os"
)

func runCommand(args []string) {
	if len(args) < 1 { // tell user how to use command
		fmt.Println("Usage: flutter-devtool <command>")
		return
	}

	command := args[0] // get first argument // for flutter-devtool devices it is devices
	switch command {
	case "devices":
	default:
		fmt.Println("Unknown commands:", command)
		fmt.Println("Available commands: devices")
	}
}

func getArguments() [] string {
	return os.Args[1:] // no need of os.args[0] so return everything after it
}