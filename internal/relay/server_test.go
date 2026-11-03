package relay

import "testing"

func TestParseHandshake(t *testing.T) {
	tests := []struct {
		name    string
		message string
		command string
		code    string
		role    ClientRole
	}{
		{
			name:    "old laptop handshake",
			message: "LAPTOP",
			command: "NEW",
			role:    Laptop,
		},
		{
			name:    "new phone handshake",
			message: "NEW PHONE",
			command: "NEW",
			role:    Phone,
		},
		{
			name:    "laptop reconnect",
			message: "RECONNECT Session-1 LAPTOP",
			command: "RECONNECT",
			code:    "Session-1",
			role:    Laptop,
		},
		{
			name:    "phone reconnect",
			message: "RECONNECT Session-1 PHONE",
			command: "RECONNECT",
			code:    "Session-1",
			role:    Phone,
		},
		{
			name:    "invalid handshake",
			message: "INVALID",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command, code, role := parseHandshake(test.message)

			if command != test.command {
				t.Errorf("Expected command %q, got %q", test.command, command)
			}

			if code != test.code {
				t.Errorf("Expected code %q, got %q", test.code, code)
			}

			if role != test.role {
				t.Errorf("Expected role %q, got %q", test.role, role)
			}
		})
	}
}