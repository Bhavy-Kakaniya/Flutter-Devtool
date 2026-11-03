package relay

import (
	"net"
	"testing"
)

func TestReconnectClient(t *testing.T) {
	manager := NewSessionManager()

	laptop, laptopPeer := net.Pipe()
	defer laptop.Close()
	defer laptopPeer.Close()

	phone, phonePeer := net.Pipe()
	defer phone.Close()
	defer phonePeer.Close()

	session := manager.AddClient(laptop, Laptop)
	manager.AddClient(phone, Phone)
	session.RemoveClient(laptop, session.generation)

	if session.Laptop != nil {
		t.Fatal("expected laptop to be disconnected")
	}

	if session.Phone != phone {
		t.Fatal("expected existing phone connection to remain")
	}

	newLaptop, newLaptopPeer := net.Pipe()
	defer newLaptop.Close()
	defer newLaptopPeer.Close()

	reconnectedSession := manager.ReconnectClient(newLaptop, session.Code, Laptop)
	if reconnectedSession == nil {
		t.Fatal("Expected laptop to reconnect, got nil")
	}
	if !session.IsReady() {
		t.Fatal("expected session to be ready after laptop reconnect")
	}
	// Disconnect the phone explicitly before reconnecting it.
	session.RemoveClient(phone, session.generation)

	// Verify the laptop remains connected.
	if session.Laptop != newLaptop {
		t.Fatal("Expected laptop to remain connected after phone disconnect")
	}

	// Verify the phone has been removed.
	if session.Phone != nil {
		t.Fatal("Expected phone to be disconnected")
	}
	if reconnectedSession.ID != session.ID {
		t.Fatalf("Expected session Id %d, got %d", session.ID, reconnectedSession.ID)
	}

	if reconnectedSession.Laptop != newLaptop {
		t.Fatal("Expected session to ready after reconnection")
	}

	currentGeneration := session.generation
	session.RemoveClient(laptop, currentGeneration-1)
	if session.Laptop != newLaptop {
		t.Fatal("Old relay removed the new laptop connection")
	}

	newPhone, newPhonePeer := net.Pipe()
	defer newPhone.Close()
	defer newPhonePeer.Close()

	reconnectedPhoneSession := manager.ReconnectClient(newPhone, session.Code, Phone)

	if reconnectedPhoneSession == nil {
		t.Fatal("Expected phone to reconnect, got nil")
	}

	if !reconnectedPhoneSession.IsReady() {
		t.Fatal("Expected session to be ready after both client reconnect")
	}
}
