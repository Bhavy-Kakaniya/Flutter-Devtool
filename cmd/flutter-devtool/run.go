package main

import (
	"devtool/internal/flutter"
	"fmt"
)

func runFlutter(args []string) {
	// expected: flutter-devtool run --device <device-id>
	// we need at least two arguments "--device" and the device ID
	if len(args) == 0 {
		fmt.Println("Usage: flutter-devtool run --device <device-id>")
		return
	}
	// first argument must be --device
	if args[0] != "--device" {
		fmt.Println("Expected --device")
		return
	}
	if len(args) < 2 {
		fmt.Println("Usage: flutter-devtool run --device <device-id>")
		return
	}

	deviceID := args[1]

	if err := flutter.Run(deviceID); err != nil {
		fmt.Println("Error:", err)
	}
}
