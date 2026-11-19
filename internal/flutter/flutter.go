package flutter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// Run starts Flutter on the selected Android device.
func Run(deviceID string, projectPath string) error {
	if projectPath == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectPath = currentDir
	}

	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	if _, err := os.Stat(pubspecPath); err != nil {
		return fmt.Errorf("no pubspec.yaml found in project directory: %s", projectPath)
	}

	command := exec.Command("flutter", "run", "-d", deviceID)
	command.Dir = projectPath
	command.Stdin = os.Stdin

	stdout, err := command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to capture Flutter output: %w", err)
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to capture Flutter errors: %w", err)
	}

	fmt.Println("Starting Flutter app on device:", deviceID)
	if err := command.Start(); err != nil {
		return fmt.Errorf("failed to start Flutter: %w", err)
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go forwardOutput(stdout, os.Stdout, &wg)
	go forwardOutput(stderr, os.Stderr, &wg)
	wg.Wait()

	return command.Wait()
}

// forwardOutput reads captured output and forwards it to the terminal.
func forwardOutput(reader io.ReadCloser, destination *os.File, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if IsAppLog(line) {
			fmt.Fprintln(destination, "[APP]", line)
		} else {
			fmt.Fprintln(destination, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Output reading error:", err)
	}
}
