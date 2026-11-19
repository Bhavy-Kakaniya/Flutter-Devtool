package hotspot

import (
	"testing"
)

func TestFindHotspotLive(t *testing.T) {
	info, err := FindHotspot()
	if err != nil {
		t.Logf("FindHotspot live returned error (may happen if offline): %v", err)
		return
	}
	t.Logf("Detected Hotspot Info:")
	t.Logf("  Gateway IP (Phone): %s", info.GatewayIP)
	t.Logf("  Laptop IP:          %s", info.LaptopIP)
	t.Logf("  Interface:          %s", info.InterfaceName)
	t.Logf("  Likely Hotspot:     %v", info.IsLikelyHotspot)

	if info.GatewayIP == "" {
		t.Error("expected non-empty GatewayIP")
	}
	if info.LaptopIP == "" {
		t.Error("expected non-empty LaptopIP")
	}
}
