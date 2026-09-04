package adb

import "testing"

func TestParseDevices(t *testing.T) {
	output := `List of devices attached
10BEAN21SX004UB	device
1382831542050146	device`
	devices := ParseDevices(output)
	if len(devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(devices))
	}

	// check the ID and status of first device
	if devices[0].ID != "10BEAN21SX004UB" {
		t.Errorf("unexpected first device ID: %s", devices[0].ID)
	}

	if devices[0].Status != "device" {
		t.Errorf("unexpected first device status: %s", devices[0].Status)
	}

	if devices[1].ID != "1382831542050146" {
		t.Errorf("unexpected second device ID: %s", devices[1].ID)
	}

	if devices[1].Status != "device" {
		t.Errorf("unexpected second device status: %s", devices[1].Status)
	}
}
