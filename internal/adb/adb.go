package adb

import (
	"fmt"
	"os/exec"
)

// GetDevices runs the ADB "devices" command
// and returns the connected Android devices.
func GetDevices(adbPath string) ([]Device, error) {

	// Create a command that will execute <adbPath> devices
	command := exec.Command(adbPath, "devices")

	// execute command and capture output
	output, err := command.Output()

	// check whether ADB command failed
	if err != nil {
		return nil, fmt.Errorf("failed to run ADB: %w", err)
	}

	// convert raw ADB text into Device structs
	devices := ParseDevices(string(output))

	return devices, nil
}