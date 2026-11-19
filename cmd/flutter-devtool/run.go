package main

import (
	"devtool/internal/adb"
	"devtool/internal/flutter"
	"fmt"
)

func runFlutter(args []string) {
	var deviceID string
	var projectPath string

	// Read the optional command-line arguments.
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--device":
			if i+1 >= len(args) {
				fmt.Println("Missing device ID")
				return
			}
			deviceID = args[i+1]
			i++

		case "--project":
			if i+1 >= len(args) {
				fmt.Println("Missing project path")
				return
			}
			projectPath = args[i+1]
			i++

		default:
			fmt.Println("Unknown argument:", args[i])
			return
		}
	}

	// If the user didn't provide a device, discover connected devices.
	if deviceID == "" {
		adbPath, err := adb.FindADB()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		devices, err := adb.GetDevices(adbPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Ask the user to select a device or select it automatically.
		deviceID, err = selectDevice(devices)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// Start Flutter on the selected device.
	if err := flutter.Run(deviceID, projectPath); err != nil {
		fmt.Println("Error:", err)
	}
}