package hotspot

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FindHotspot automatically detects the phone's gateway IP and laptop's IP over the hotspot.
func FindHotspot() (*HotspotInfo, error) {
	// First attempt: parse ipconfig (accurate adapter naming on Windows)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ipconfig")
	output, err := cmd.Output()
	if err == nil {
		adapters, parseErr := ParseIPConfig(string(output))
		if parseErr == nil && len(adapters) > 0 {
			// Prefer wireless adapter marked as likely hotspot
			for _, adp := range adapters {
				if adp.IsLikelyHotspot {
					return &adp, nil
				}
			}
			// Fallback to the first available adapter with gateway
			return &adapters[0], nil
		}
	}

	// Second attempt: parse route print 0.0.0.0
	ctxRoute, cancelRoute := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelRoute()

	routeCmd := exec.CommandContext(ctxRoute, "route", "print", "0.0.0.0")
	routeOutput, routeErr := routeCmd.Output()
	if routeErr == nil {
		gw, iface, err := ParseRoutePrint(string(routeOutput))
		if err == nil {
			return &HotspotInfo{
				GatewayIP:       gw,
				LaptopIP:        iface,
				InterfaceName:   "Default Route Interface",
				IsLikelyHotspot: isHotspotSubnet(gw),
			}, nil
		}
	}

	return nil, fmt.Errorf("unable to detect phone hotspot gateway IP")
}

// CheckPortReachable tests whether a specific TCP port on the target host is reachable.
func CheckPortReachable(host string, port int, timeout time.Duration) bool {
	target := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// CheckADB tests whether ADB (port 5555) is listening on the phone gateway.
func CheckADB(gatewayIP string) bool {
	return CheckPortReachable(gatewayIP, 5555, 1500*time.Millisecond)
}

// ConnectADB attempts to connect ADB to the phone's gateway IP on port 5555.
func ConnectADB(adbPath string, gatewayIP string) (string, error) {
	return ConnectADBPort(adbPath, gatewayIP, 5555)
}

// ConnectADBPort connects ADB to a specific port on the target IP.
func ConnectADBPort(adbPath string, ip string, port int) (string, error) {
	target := fmt.Sprintf("%s:%d", ip, port)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, adbPath, "connect", target)
	output, err := cmd.CombinedOutput()
	result := strings.TrimSpace(string(output))
	if err != nil {
		return result, fmt.Errorf("adb connect failed: %w (output: %s)", err, result)
	}
	return result, nil
}

// PairADB executes `adb pair <ip>:<port> <code>` using Android 11+ native pairing.
func PairADB(adbPath string, ip string, port string, code string) (string, error) {
	target := fmt.Sprintf("%s:%s", ip, port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, adbPath, "pair", target, code)
	output, err := cmd.CombinedOutput()
	result := strings.TrimSpace(string(output))
	if err != nil || strings.Contains(strings.ToLower(result), "failed") {
		return result, fmt.Errorf("pairing failed: %s", result)
	}
	return result, nil
}

// FindOpenWirelessPorts rapidly scans the common Android wireless debugging port range (30000 - 48000).
func FindOpenWirelessPorts(ip string) []int {
	var openPorts []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, 250) // 250 parallel probes

	for port := 30000; port <= 48000; port++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()

			conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(p)), 200*time.Millisecond)
			if err == nil {
				conn.Close()
				mu.Lock()
				openPorts = append(openPorts, p)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return openPorts
}

// SetTCPIPPort sets the ADB daemon on the device to also listen on port 5555.
func SetTCPIPPort(adbPath string, targetDevice string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, adbPath, "-s", targetDevice, "tcpip", "5555")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set adb tcpip: %w (output: %s)", err, string(output))
	}
	return nil
}
