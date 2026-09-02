package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Device struct {
	ID     string
	Status string
}

func findADB() (string, error) {
	// 1. Check if ADB is available in PATH.
	if adbPath, err := exec.LookPath("adb"); err == nil {
		return adbPath, nil
	}

	// 2. Check ANDROID_HOME.
	if androidHome := os.Getenv("ANDROID_HOME"); androidHome != "" {
		if adbPath := adbPathFromSDK(androidHome); adbPath != "" {
			return adbPath, nil
		}
	}

	// 3. Check ANDROID_SDK_ROOT.
	if androidSDKRoot := os.Getenv("ANDROID_SDK_ROOT"); androidSDKRoot != "" {
		if adbPath := adbPathFromSDK(androidSDKRoot); adbPath != "" {
			return adbPath, nil
		}
	}

	// 4. Ask Flutter for its configured Android SDK.
	if sdkPath, err := findFlutterAndroidSDK(); err == nil {
		if adbPath := adbPathFromSDK(sdkPath); adbPath != "" {
			return adbPath, nil
		}
	}

	// 5. Check the standard Windows Android SDK location.
	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		sdkPath := filepath.Join(localAppData, "Android", "Sdk")
		if adbPath := adbPathFromSDK(sdkPath); adbPath != "" {
			return adbPath, nil
		}
	}
	return "", fmt.Errorf("ADB executable not found")
}

func adbPathFromSDK(sdkPath string) string {
	adbPath := filepath.Join(sdkPath, "platform-tools", "adb.exe")
	if _, err := os.Stat(adbPath); err == nil {
		return adbPath
	}
	return ""
}

func findFlutterAndroidSDK() (string, error) {
	command := exec.Command("flutter", "config", "--list")
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "android-sdk:") {
			sdkPath := strings.TrimSpace(
				strings.TrimPrefix(line, "android-sdk:"),
			)
			return sdkPath, nil
		}
	}
	return "", fmt.Errorf("Android SDK not configured in Flutter")
}

func parseDevices(output string) []Device {
	var devices []Device

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "List of devices attached") {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		device := Device{
			ID:     fields[0],
			Status: fields[1],
		}

		devices = append(devices, device)
	}

	return devices
}

func getDevices(adbPath string) ([]Device, error) {
	command := exec.Command(adbPath, "devices") // command "adb devices"

	output, err := command.Output() // run command and get what it prints
	if err != nil {
		return nil, err // if err return no device
	}
	devices := parseDevices(string(output)) // convert raw adb output in Device struct
	return devices, nil
}

func runDevices() {
	// find where ADB is installed on computer
	adbPath, err := findADB()

	// stop id ADB not found
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// show adb location so user know which adb is being used
	fmt.Println("ADB found at", adbPath)

	// ask ADB for all currently connected devices
	devices, err := getDevices(adbPath)

	// stop if ADB failed to return device list
	if err != nil {
		fmt.Println("Error running ADB:", err)
		return
	}

	fmt.Println("\nConneced devices")
	if len(devices) == 0 {
		// tell user no device were found
		fmt.Println("No device connected")
		return
	}

	// loop through every connected device
	for _, device := range devices {
		fmt.Println("ID:", device.ID)
		fmt.Println("Status:", device.Status)
	}
}

func main() {

	//check whether user provided a command after program name like go run . devices
	if len(os.Args) < 2 {
		// Tell user how program should be used
		fmt.Println("Usage: flutter-devtool <command>")

		// stop program as no command was provided
		return
	}

	// get first argument provided by user
	// example: go run . devices, devices will be first argument
	command := os.Args[1]

	// check which command is requested
	switch command {
	case "devices":
		// run the devices command
		runDevices()
	default:
		// user entered command which is not supported
		fmt.Println("Unknown command:", command)
		fmt.Println("Available commands: devices")
	}
}