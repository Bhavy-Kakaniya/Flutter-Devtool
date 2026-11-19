package main

import (
	"fmt"
	"os"
)

// getArguments returns the arguments provided by the user
// os.Args[0] is the program itself, so we remove it
func getArguments() []string {
	return os.Args[1:]
}

// runCommand decides which command the user requested
func runCommand(args []string) {

	// check whether the user provided any command
	if len(args) < 1 {
		fmt.Println("Usage: flutter-devtool <command>")
		fmt.Println("\nAvailable commands:")
		fmt.Println("  hotspot  - Check phone hotspot connection and status")
		fmt.Println("  pair     - Pair with phone using Android native Wireless Debugging (No 3rd party apps)")
		fmt.Println("  connect  - Automatically discover phone over hotspot and connect ADB")
		fmt.Println("  devices  - List ADB connected devices")
		fmt.Println("  run      - Run Flutter app on selected device")
		fmt.Println("  relay    - Start TCP relay server for companion agent")
		return
	}

	// get the first user-provided argument example: ["devices"] -> "devices"
	command := args[0]

	// decide what to execute based on command
	switch command {
	case "hotspot":
		runHotspot()
	case "pair":
		runPair(args[1:])
	case "connect":
		runConnect()
	case "devices":
		runDevices()
	case "run":
		runFlutter(args[1:])
	case "relay":
		runRelay()
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Available commands:")
		fmt.Println("  hotspot  - Check phone hotspot connection and status")
		fmt.Println("  pair     - Pair with phone using Android native Wireless Debugging (No 3rd party apps)")
		fmt.Println("  connect  - Automatically discover phone over hotspot and connect ADB")
		fmt.Println("  devices  - List ADB connected devices")
		fmt.Println("  run      - Run Flutter app on selected device")
		fmt.Println("  relay    - Start TCP relay server for companion agent")
	}
}
