package adb

import "strings"

type Device struct {
	ID     string
	Status string
}

// ParseDevices converts the raw output from "adb devices" into slice of Device structs
func ParseDevices(output string) []Device {
	var devices []Device
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "List of devices attached") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		device := Device{
			ID:     fields[0],
			Status: fields[1],
		}

		devices = append(devices, device)
	}

	return devices
}
