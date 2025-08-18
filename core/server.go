package core

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// sendServerMessage sends a text message to the server
func sendServerMessage(server, port, encryption, message string, silent bool) error {
	protocol := "https"
	if encryption == "no" {
		protocol = "http"
	}
	url := fmt.Sprintf("%s://%s:%s/message", protocol, server, port)
	data := strings.NewReader(fmt.Sprintf("message=%s", message))
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to create server message request: %v%s\n",
				ColorYellow, ColorWhite, err, ColorReset)
		}
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to send server message: %v%s\n",
				ColorYellow, ColorWhite, err, ColorReset)
		}
		return fmt.Errorf("failed to send server message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if !silent {
			fmt.Printf("%s[WARNING]%s Server message failed with status %d: %s%s\n",
				ColorYellow, ColorWhite, resp.StatusCode, string(bodyBytes), ColorReset)
		}
		return fmt.Errorf("server message failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if !silent {
		fmt.Printf("%s[INFO]%s Message sent to server: %s%s\n",
			ColorGreen, ColorWhite, message, ColorReset)
	}
	return nil
}

// SendToServer sends a file to the specified server
func SendToServer(server, port, encryption, filePath string, silent bool) error {
	// Validate file existence
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filePath)
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to create form file for %s: %v%s\n",
				ColorYellow, ColorWhite, filePath, err, ColorReset)
		}
		return fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to write file %s to form: %v%s\n",
				ColorYellow, ColorWhite, filePath, err, ColorReset)
		}
		return fmt.Errorf("failed to write file to form: %v", err)
	}
	if err := writer.Close(); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to close form data for %s: %v%s\n",
				ColorYellow, ColorWhite, filePath, err, ColorReset)
		}
		return fmt.Errorf("failed to close form data: %v", err)
	}

	// Create and send request
	protocol := "https"
	if encryption == "no" {
		protocol = "http"
	}
	url := fmt.Sprintf("%s://%s:%s/upload", protocol, server, port)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to create server request for %s: %v%s\n",
				ColorYellow, ColorWhite, filePath, err, ColorReset)
		}
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to send file %s to server: %v%s\n",
				ColorYellow, ColorWhite, filePath, err, ColorReset)
		}
		return fmt.Errorf("failed to send file to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if !silent {
			fmt.Printf("%s[WARNING]%s Server file upload failed with status %d: %s%s\n",
				ColorYellow, ColorWhite, resp.StatusCode, string(bodyBytes), ColorReset)
		}
		return fmt.Errorf("server file upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if !silent {
		fmt.Printf("%s[INFO]%s File %s sent to server%s\n",
			ColorGreen, ColorWhite, filepath.Base(filePath), ColorReset)
	}
	return nil
}
