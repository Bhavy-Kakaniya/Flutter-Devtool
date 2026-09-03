package adb

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

func FindADB() (string, error) {
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