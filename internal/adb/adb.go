package adb

import (
	"fmt"
	"os/exec"
)

// GetDevices runs the ADB "devices" command and returns the connected Android devices.
func GetDevices(adbPath string) ([]Device, error) {
	command := exec.Command(adbPath, "devices")
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run ADB: %w", err)
	}
	devices := ParseDevices(string(output))
	return devices, nil
}