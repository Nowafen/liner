package payloads

import (
	"runtime"
)

// DetectOS returns the OS information
func DetectOS() (string, error) {
	return runtime.GOOS, nil
}
