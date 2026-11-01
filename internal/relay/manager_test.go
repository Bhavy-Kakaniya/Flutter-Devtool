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
	session.RemoveClient(laptop)

	newLaptop, newLaptopPeer :=  net.Pipe()
	defer newLaptop.Close()
	defer newLaptopPeer.Close()

	reconnectedSession := manager.RecconectClient(newLaptop, session.Code, Laptop)
	if reconnectedSession == nil{
		t.Fatal("Expected laptop to reconnect, got nil")
	}

	if reconnectedSession.ID != session.ID{
		t.Fatalf("Expected session Id %d, got %d", session.ID, reconnectedSession.ID)
	}

	if reconnectedSession.Laptop != newLaptop {
		t.Fatal("Expected session to ready after reconnection")
	}
	if reconnectedSession.Phone != nil {
		t.Fatal("Expected phone connection to be nil")
	}

	newPhone, newPhonePeer := net.Pipe()
	defer newPhone.Close()
	defer newPhonePeer.Close()

	reconnectedPhoneSession := manager.RecconectClient(newPhone, session.Code, Phone)

	if reconnectedPhoneSession == nil {
		t.Fatal("Expected phone to reconnect, got nil")
	}
	if !reconnectedPhoneSession.IsReady() {
		t.Fatal("Expected session to be ready after both client reconnect")
	}
}
