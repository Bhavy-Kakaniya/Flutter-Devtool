package flutter

import "testing"

func TestIsAppLog(t *testing.T) {
	appLogs := []string{
		"I/flutter (15118): Hello Flutter",
		"E/flutter (15118): Something went wrong",
		"W/flutter (15118): Warning",
		"D/flutter (15118): Debug message",
	}

	for _, line := range appLogs {
		if !IsAppLog(line) {
			t.Errorf("expected application log: %s", line)
		}
	}

	otherLogs := []string{
		"Running Gradle task assembleDebug...",
		"Installing build/app-debug.apk...",
		"Syncing files to device...",
	}

	for _, line := range otherLogs {
		if IsAppLog(line) {
			t.Errorf("unexpected application log: %s", line)
		}
	}
}
