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

	_, err := phonePeer.Write([]byte("test"))
	if err == nil {
		t.Fatal("Expected old phone connection to be closed")
	}

	newLaptop, newLaptopPeer :=  net.Pipe()
	defer newLaptop.Close()
	defer newLaptopPeer.Close()

	reconnectedSession := manager.ReconnectClient(newLaptop, session.Code, Laptop)
	if reconnectedSession == nil{
		t.Fatal("Expected laptop to reconnect, got nil")
	}

	if reconnectedSession.ID != session.ID{
		t.Fatalf("Expected session Id %d, got %d", session.ID, reconnectedSession.ID)
	}

	if reconnectedSession.Laptop != newLaptop {
		t.Fatal("Expected session to ready after reconnection")
	}

	currentGeneration := session.generation
	session.RemoveClient(laptop, currentGeneration - 1)
	if session.Laptop != newLaptop {
		t.Fatal("Old relay removed the new laptop connection")
	}

	if reconnectedSession.Phone != nil {
		t.Fatal("Expected phone connection to be nil")
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
