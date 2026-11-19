package flutter

import "strings"

func IsAppLog(line string) bool {
	line = strings.TrimSpace(line)
	return strings.Contains(line, "I/flutter") ||
		strings.Contains(line, "E/flutter") ||
		strings.Contains(line, "W/flutter") ||
		strings.Contains(line, "D/flutter")
}
