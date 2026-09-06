package flutter

import (
	"fmt"
	"os"
	"os/exec"
)

// run starts "flutter run" for given device
// deviceID tells flutter which Android device should receive the application

func Run(deviceID string) error {
	command := exec.Command("flutter", "run", "-d", deviceID)

	// connect flutter's standard input, outoput, error to our terminal this allows user to interact with "flutter run"
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	fmt.Println("Starting flutter app on device:", deviceID)
	return command.Run() // start flutter and wait until it finishes
}
