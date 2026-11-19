package hotspot

import (
	"testing"
)

func TestParseRoutePrint(t *testing.T) {
	output := `===========================================================================
Interface List
 10...f8 54 f6 60 14 e7 ......MediaTek Wi-Fi 6 MT7921 Wireless LAN Card
===========================================================================

IPv4 Route Table
===========================================================================
Active Routes:
Network Destination        Netmask          Gateway       Interface  Metric
          0.0.0.0          0.0.0.0     10.235.10.94    10.235.10.157     55
===========================================================================
`
	gw, iface, err := ParseRoutePrint(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gw != "10.235.10.94" {
		t.Errorf("expected gateway 10.235.10.94, got %s", gw)
	}
	if iface != "10.235.10.157" {
		t.Errorf("expected interface 10.235.10.157, got %s", iface)
	}
}

func TestParseIPConfig(t *testing.T) {
	output := `Windows IP Configuration

Ethernet adapter Ethernet:

   Media State . . . . . . . . . . . : Media disconnected
   Connection-specific DNS Suffix  . : 

Wireless LAN adapter Wi-Fi:

   Connection-specific DNS Suffix  . : 
   IPv4 Address. . . . . . . . . . . : 10.235.157.157
   Subnet Mask . . . . . . . . . . . : 255.255.255.0
   Default Gateway . . . . . . . . . : 10.235.157.208
`
	results, err := ParseIPConfig(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 adapter result, got %d", len(results))
	}
	info := results[0]
	if info.GatewayIP != "10.235.157.208" {
		t.Errorf("expected gateway 10.235.157.208, got %s", info.GatewayIP)
	}
	if info.LaptopIP != "10.235.157.157" {
		t.Errorf("expected laptop IP 10.235.157.157, got %s", info.LaptopIP)
	}
	if !info.IsLikelyHotspot {
		t.Errorf("expected IsLikelyHotspot to be true")
	}
}

func TestIsHotspotSubnet(t *testing.T) {
	cases := []struct {
		ip       string
		expected bool
	}{
		{"192.168.43.1", true},
		{"172.20.10.1", true},
		{"10.235.157.208", true},
		{"192.168.1.1", false},
		{"127.0.0.1", false},
	}

	for _, c := range cases {
		actual := isHotspotSubnet(c.ip)
		if actual != c.expected {
			t.Errorf("isHotspotSubnet(%s): expected %v, got %v", c.ip, c.expected, actual)
		}
	}
}
