package main

import (
	"fmt"

	"devtool/internal/adb"
)

// runDevices handles the "devices" CLI command.
func runDevices() {

	// Find the ADB executable.
	adbPath, err := adb.FindADB()

	// Stop if ADB cannot be found.
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Show which ADB executable is being used.
	fmt.Println("ADB found at:", adbPath)

	// Ask the ADB package for connected devices.
	devices, err := adb.GetDevices(adbPath)

	// Stop if ADB fails.
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Display a heading.
	fmt.Println("\nConnected devices")

	// Handle the case where no devices are connected.
	if len(devices) == 0 {
		fmt.Println("No devices connected.")
		return
	}

	// Display every connected device.
	for _, device := range devices {
		fmt.Println("ID:", device.ID)
		fmt.Println("Status:", device.Status)
	}
}