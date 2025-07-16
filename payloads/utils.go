package payloads

import (
	"os/exec"
	"runtime"
	"strings"
)

// DetectOS returns the OS information using system commands and runtime
func DetectOS() (string, error) {
	// Try using uname command
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err == nil {
		osName := strings.TrimSpace(string(output))
		return osName, nil
	}

	// Fallback to runtime.GOOS
	return runtime.GOOS, nil
}
