package hotspot

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// HotspotInfo holds discovered network information about the phone hotspot.
type HotspotInfo struct {
	GatewayIP       string `json:"gateway_ip"`       // Phone's IP address (default gateway) or empty if laptop is host
	LaptopIP        string `json:"laptop_ip"`        // Laptop's assigned IP address
	SubnetMask      string `json:"subnet_mask"`      // Subnet mask
	InterfaceName   string `json:"interface_name"`   // e.g. "Wi-Fi" or "Wireless LAN adapter Wi-Fi"
	IsLikelyHotspot bool   `json:"is_likely_hotspot"` // True if gateway matches common hotspot subnets
	IsLaptopHotspot bool   `json:"is_laptop_hotspot"` // True if laptop is hosting the hotspot (192.168.137.x)
}

// ParseRoutePrint parses the output of `route print 0.0.0.0` to extract the default gateway and interface IP.
func ParseRoutePrint(output string) (gatewayIP string, interfaceIP string, err error) {
	lines := strings.Split(output, "\n")
	foundActiveRoutes := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Active Routes:") {
			foundActiveRoutes = true
			continue
		}
		if !foundActiveRoutes {
			continue
		}
		if strings.HasPrefix(line, "==") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		// Expected format: Network Destination, Netmask, Gateway, Interface, Metric
		// e.g.: 0.0.0.0 0.0.0.0 10.235.10.94 10.235.10.157 55
		if len(fields) >= 5 && fields[0] == "0.0.0.0" && fields[1] == "0.0.0.0" {
			gw := fields[2]
			iface := fields[3]
			if net.ParseIP(gw) != nil && net.ParseIP(iface) != nil && gw != "0.0.0.0" {
				return gw, iface, nil
			}
		}
	}
	return "", "", fmt.Errorf("no active default route found in route print output")
}

// ParseIPConfig parses the output of `ipconfig` looking specifically for wireless adapters with a default gateway.
func ParseIPConfig(output string) ([]HotspotInfo, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var results []HotspotInfo
	var currentAdapter string
	var currentIP, currentMask, currentGW string

	flushCurrent := func() {
		if currentAdapter != "" && currentIP != "" {
			if currentGW != "" && currentGW != "0.0.0.0" {
				isWireless := strings.Contains(strings.ToLower(currentAdapter), "wireless") ||
					strings.Contains(strings.ToLower(currentAdapter), "wi-fi") ||
					strings.Contains(strings.ToLower(currentAdapter), "wlan")

				info := HotspotInfo{
					GatewayIP:       currentGW,
					LaptopIP:        currentIP,
					SubnetMask:      currentMask,
					InterfaceName:   currentAdapter,
					IsLikelyHotspot: isWireless || isHotspotSubnet(currentGW),
					IsLaptopHotspot: false,
				}
				results = append(results, info)
			} else if strings.HasPrefix(currentIP, "192.168.137.") {
				// Windows Mobile Hotspot adapter
				info := HotspotInfo{
					GatewayIP:       "",
					LaptopIP:        currentIP,
					SubnetMask:      currentMask,
					InterfaceName:   currentAdapter + " (Windows Mobile Hotspot)",
					IsLikelyHotspot: true,
					IsLaptopHotspot: true,
				}
				results = append(results, info)
			}
		}
		currentIP = ""
		currentMask = ""
		currentGW = ""
	}

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Adapter header line (starts at beginning of line, not indented, ends with ':')
		if len(line) > 0 && line[0] != ' ' && line[0] != '\t' && strings.HasSuffix(trimmed, ":") {
			flushCurrent()
			currentAdapter = strings.TrimSuffix(trimmed, ":")
			continue
		}

		if strings.Contains(trimmed, "IPv4 Address") {
			parts := strings.Split(trimmed, ":")
			if len(parts) >= 2 {
				currentIP = strings.TrimSpace(parts[len(parts)-1])
				// Clean any trailing info like (Preferred)
				currentIP = strings.Fields(currentIP)[0]
			}
		} else if strings.Contains(trimmed, "Subnet Mask") {
			parts := strings.Split(trimmed, ":")
			if len(parts) >= 2 {
				currentMask = strings.TrimSpace(parts[len(parts)-1])
			}
		} else if strings.Contains(trimmed, "Default Gateway") {
			parts := strings.Split(trimmed, ":")
			if len(parts) >= 2 {
				val := strings.TrimSpace(parts[len(parts)-1])
				if val != "" && net.ParseIP(val) != nil {
					currentGW = val
				}
			}
		}
	}
	flushCurrent()

	if len(results) == 0 {
		return nil, fmt.Errorf("no active adapter with default gateway found in ipconfig")
	}
	return results, nil
}

// isHotspotSubnet checks if an IP belongs to standard mobile hotspot ranges:
// - Android typical: 192.168.43.1, 192.168.50.1, 10.x.x.x
// - iOS typical: 172.20.10.1
func isHotspotSubnet(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	// 192.168.43.0/24 is standard Android hotspot
	if ip4[0] == 192 && ip4[1] == 168 && (ip4[2] == 43 || ip4[2] == 50) {
		return true
	}
	// 172.20.10.0/24 is standard iOS hotspot
	if ip4[0] == 172 && ip4[1] == 20 && ip4[2] == 10 {
		return true
	}
	// 10.x.x.x is common on some carriers / modern Android SoftAP
	if ip4[0] == 10 {
		return true
	}
	return false
}
