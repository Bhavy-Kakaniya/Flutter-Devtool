package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"devtool/internal/adb"
	"devtool/internal/hotspot"
)

// runHotspot displays current hotspot connection details.
func runHotspot() {
	info, err := hotspot.FindHotspot()
	if err != nil {
		fmt.Println("No active hotspot detected:", err)
		return
	}

	fmt.Println("--- Hotspot Network Information ---")
	fmt.Printf("  Interface:       %s\n", info.InterfaceName)
	if info.IsLaptopHotspot {
		fmt.Println("  Mode:            Laptop is broadcasting Windows Mobile Hotspot")
		fmt.Printf("  Laptop Host IP:  %s\n", info.LaptopIP)
		fmt.Println("  Phone IP:        Check on phone Settings -> Wi-Fi -> Network details (192.168.137.x)")
		return
	}
	fmt.Printf("  Phone IP (GW):   %s\n", info.GatewayIP)
	fmt.Printf("  Laptop IP:       %s\n", info.LaptopIP)
	if info.SubnetMask != "" {
		fmt.Printf("  Subnet Mask:     %s\n", info.SubnetMask)
	}

	adbOpen := hotspot.CheckADB(info.GatewayIP)
	if adbOpen {
		fmt.Println("  Wireless ADB:    AVAILABLE (Port 5555 open)")
	} else {
		fmt.Println("  Wireless ADB:    NOT ACTIVE (Port 5555 closed)")
	}
}

// runPair pairs with Android 11+ native Wireless Debugging without any 3rd party apps.
func runPair(args []string) {
	fmt.Println("=== Android Native Wireless Pairing ===")
	info, _ := hotspot.FindHotspot()

	adbPath, err := adb.FindADB()
	if err != nil {
		fmt.Println("Error locating ADB:", err)
		return
	}

	defaultIP := ""
	if info != nil {
		defaultIP = info.GatewayIP
	}

	fmt.Println("Instructions on your Android phone (Android 11+):")
	fmt.Println("  1. Open Settings -> Developer options -> Wireless debugging.")
	fmt.Println("  2. Turn Wireless debugging ON.")
	fmt.Println("  3. Tap 'Pair device with pairing code'.")
	fmt.Println("  4. A popup will show a 6-digit Wi-Fi pairing code and an IP:Port.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Ask for target IP (with default)
	targetIP := defaultIP
	if defaultIP != "" {
		fmt.Printf("Phone IP [%s] (press Enter to accept or enter IP shown on phone): ", defaultIP)
	} else {
		fmt.Print("Enter Phone IP shown on the pairing popup: ")
	}
	ipInput, _ := reader.ReadString('\n')
	trimmedIP := strings.TrimSpace(ipInput)
	if trimmedIP != "" {
		if strings.Contains(trimmedIP, ":") {
			parts := strings.Split(trimmedIP, ":")
			targetIP = parts[0]
		} else {
			targetIP = trimmedIP
		}
	}

	if targetIP == "" {
		fmt.Println("Phone IP cannot be empty.")
		return
	}

	if strings.HasPrefix(targetIP, "192.0.0.") {
		fmt.Printf("\n[!] %s is your phone's cellular mobile data (464XLAT) IP, NOT a Wi-Fi IP!\n", targetIP)
		fmt.Println("Your laptop cannot connect to the phone over cellular data.")
		fmt.Println("\nTo fix this:")
		fmt.Println("  1. Turn OFF Mobile Data on your phone.")
		fmt.Println("  2. Check your phone's real Wi-Fi IP (Settings -> Wi-Fi -> tap your network -> Network details).")
		fmt.Println("  3. Tap 'Pair device with pairing code' again on your phone and re-run this command.")
		return
	}

	// Ask for pairing port
	fmt.Print("Enter the pairing port from the popup (e.g. 38491): ")
	portInput, _ := reader.ReadString('\n')
	pairingPort := strings.TrimSpace(portInput)
	if pairingPort == "" {
		fmt.Println("Pairing port cannot be empty.")
		return
	}

	// Ask for 6-digit pairing code
	fmt.Print("Enter the 6-digit Wi-Fi pairing code: ")
	codeInput, _ := reader.ReadString('\n')
	pairingCode := strings.TrimSpace(codeInput)
	if pairingCode == "" {
		fmt.Println("Pairing code cannot be empty.")
		return
	}

	fmt.Printf("\nPairing with %s:%s...\n", targetIP, pairingPort)
	pairOutput, err := hotspot.PairADB(adbPath, targetIP, pairingPort, pairingCode)
	if err != nil {
		fmt.Println("Pairing failed:", err)
		return
	}
	fmt.Println(pairOutput)
	fmt.Println("\nPairing successful!")

	// Look for the connection port (shown on the main Wireless debugging screen)
	fmt.Println("\nNow looking for the connection port (shown under 'IP address & Port' on the Wireless debugging screen)...")
	fmt.Print("Enter connection port (leave blank to auto-scan): ")
	connPortInput, _ := reader.ReadString('\n')
	connPortStr := strings.TrimSpace(connPortInput)

	var connectPort int
	if connPortStr != "" {
		p, err := strconv.Atoi(connPortStr)
		if err == nil {
			connectPort = p
		}
	}

	if connectPort == 0 {
		fmt.Println("Scanning phone for active wireless debugging port...")
		ports := hotspot.FindOpenWirelessPorts(targetIP)
		for _, p := range ports {
			if strconv.Itoa(p) != pairingPort {
				connectPort = p
				break
			}
		}
		if connectPort == 0 && len(ports) > 0 {
			connectPort = ports[0]
		}
	}

	if connectPort == 0 {
		fmt.Print("Could not auto-detect connection port. Please enter the port from your phone screen: ")
		manualPort, _ := reader.ReadString('\n')
		connectPort, _ = strconv.Atoi(strings.TrimSpace(manualPort))
	}

	if connectPort > 0 {
		fmt.Printf("Connecting to %s:%d...\n", targetIP, connectPort)
		connectOutput, err := hotspot.ConnectADBPort(adbPath, targetIP, connectPort)
		if err != nil {
			fmt.Println("Connection failed:", err)
			return
		}
		fmt.Println(connectOutput)

		// Set fixed port 5555 so future connects don't require port lookup
		target := fmt.Sprintf("%s:%d", targetIP, connectPort)
		if err := hotspot.SetTCPIPPort(adbPath, target); err == nil {
			fmt.Println("Enabled persistent ADB port 5555 on device.")
		}

		fmt.Println("\nSuccess! Your phone is now paired and connected wirelessly.")
		fmt.Println("You can now run 'flutter-devtool run' to launch your Flutter app.")
	}
}

// runConnect automatically discovers the phone gateway and connects ADB over Wi-Fi.
func runConnect() {
	fmt.Println("Searching for phone over hotspot...")
	info, err := hotspot.FindHotspot()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Found phone at %s (Laptop IP: %s)\n", info.GatewayIP, info.LaptopIP)

	adbPath, err := adb.FindADB()
	if err != nil {
		fmt.Println("Error locating ADB:", err)
		return
	}

	// 1. Try standard port 5555 first
	if hotspot.CheckADB(info.GatewayIP) {
		fmt.Printf("Connecting ADB to %s:5555...\n", info.GatewayIP)
		output, err := hotspot.ConnectADB(adbPath, info.GatewayIP)
		if err == nil {
			fmt.Println("ADB Output:", output)
			fmt.Println("\nSuccess! You can now run 'flutter-devtool run' to launch your Flutter app.")
			return
		}
	}

	// 2. Scan for open dynamic wireless debugging ports
	fmt.Println("Probing phone for active Wireless Debugging ports...")
	ports := hotspot.FindOpenWirelessPorts(info.GatewayIP)
	if len(ports) > 0 {
		for _, port := range ports {
			fmt.Printf("Attempting connection on port %d...\n", port)
			output, err := hotspot.ConnectADBPort(adbPath, info.GatewayIP, port)
			if err == nil && strings.Contains(strings.ToLower(output), "connected") {
				fmt.Println("Connected successfully:", output)
				target := fmt.Sprintf("%s:%d", info.GatewayIP, port)
				_ = hotspot.SetTCPIPPort(adbPath, target)
				fmt.Println("\nSuccess! You can now run 'flutter-devtool run' to launch your Flutter app.")
				return
			}
		}
	}

	// 3. Not connected yet - guide the user to native pair
	fmt.Printf("\n[!] Wireless Debugging is not yet active on %s.\n", info.GatewayIP)
	fmt.Println("\nTo connect natively using built-in Android features (no 3rd-party software needed):")
	fmt.Println("  1. On your phone, open Settings -> Developer options -> Wireless debugging.")
	fmt.Println("  2. Turn Wireless debugging ON and tap 'Pair device with pairing code'.")
	fmt.Println("  3. Run: flutter-devtool pair")
}
