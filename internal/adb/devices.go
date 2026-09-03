package adb

import "strings"

// Device represents one Android device detected by ADB.
type Device struct {
	ID     string
	Status string // current ADB connection state like "device", "offline", "unauthorized".
}

// ParseDevices converts the raw output from "adb devices" into slice of Device structs
func ParseDevices(output string) []Device {

	var devices []Device // store every device found in the ADB output

	lines := strings.Split(output, "\n") // split complete ADB output into individual lines

	// process every line one at a time
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// ignore header printed by ADB Example: "List of devices attached"
		if strings.HasPrefix(line, "List of devices attached") {
			continue
		}

		fields := strings.Fields(line) // split line into separate fields

		if len(fields) < 2 { // valid device line must contain atleast id and status
			continue
		}

		// Device struct using parsed fields.
		device := Device{
			ID:     fields[0],
			Status: fields[1],
		}

		// add device to slice
		devices = append(devices, device)
	}

	return devices
}
