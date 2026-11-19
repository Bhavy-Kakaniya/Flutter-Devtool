package main

import (
	"bufio"
	"devtool/internal/adb"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func selectDevice(devices []adb.Device) (string, error) {

	// Keep only devices that are authorized and ready to use.
	var availableDevices []adb.Device

	for _, device := range devices {
		if device.Status == "device" {
			availableDevices = append(availableDevices, device)
		}
	}

	// Stop if no usable Android devices are available.
	if len(availableDevices) == 0 {
		return "", fmt.Errorf("no authorized Android devices found")
	}

	// Automatically select the device if only one is available.
	if len(availableDevices) == 1 {
		fmt.Println("Automatically selected device:", availableDevices[0].ID)
		return availableDevices[0].ID, nil
	}

	// Display only usable devices.
	fmt.Println("Multiple devices found:")

	for i, device := range availableDevices {
		fmt.Printf("%d. %s\n", i+1, device.ID)
	}

	// Read the user's selection from the terminal.
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Select device number: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		// Convert the input into an integer.
		choice, err := strconv.Atoi(strings.TrimSpace(input))

		// Reject invalid numbers and selections outside the list.
		if err != nil || choice < 1 || choice > len(availableDevices) {
			fmt.Println("Invalid selection. Try again.")
			continue
		}

		// Return the selected device's ID.
		return availableDevices[choice-1].ID, nil
	}
}
