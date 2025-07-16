package core

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// sendTelegramMessage sends a text message to Telegram using application/x-www-form-urlencoded
func sendTelegramMessage(token, chatID, message string, silent bool) error {
	// Prepare form data
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		if !silent {
			fmt.Printf("[INFO] Warning: failed to create Telegram message request: %v\n", err)
		}
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 30 * time.Second, // Set timeout to avoid hanging
	}
	resp, err := client.Do(req)
	if err != nil {
		if !silent {
			fmt.Printf("[INFO] Warning: failed to send Telegram message: %v\n", err)
		}
		return fmt.Errorf("failed to send Telegram message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if !silent {
			fmt.Printf("[INFO] Warning: Telegram message failed with status %d: %s\n", resp.StatusCode, string(bodyBytes))
		}
		return fmt.Errorf("Telegram message failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if !silent {
		fmt.Printf("[INFO] Message sent to Telegram: %s\n", message)
	}
	return nil
}

// SendToTelegram sends a file to the specified Telegram chat using multipart/form-data
func SendToTelegram(token, chatID, filePath string, silent bool) error {
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
	part, err := writer.CreateFormFile("document", filepath.Base(filePath))
	if err != nil {
		if !silent {
			fmt.Printf("[INFO] Warning: failed to create form file for %s: %v\n", filePath, err)
		}
		return fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		if !silent {
			fmt.Printf("[INFO] Warning: failed to write file %s to form: %v\n", filePath, err)
		}
		return fmt.Errorf("failed to write file to form: %v", err)
	}
	_ = writer.WriteField("chat_id", chatID)
	if err := writer.Close(); err != nil {
		if !silent {
			fmt.Printf("[INFO] Warning: failed to close form data for %s: %v\n", filePath, err)
		}
		return fmt.Errorf("failed to close form data: %v", err)
	}

	// Create and send request
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", token)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		if !silent {
			fmt.Printf("[INFO] Warning: failed to create Telegram request for %s: %v\n", filePath, err)
		}
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 60 * time.Second, // Increased timeout for file uploads
	}
	resp, err := client.Do(req)
	if err != nil {
		if !silent {
			fmt.Printf("[INFO] Warning: failed to send file %s to Telegram: %v\n", filePath, err)
		}
		return fmt.Errorf("failed to send file to Telegram: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if !silent {
			fmt.Printf("[INFO] Warning: Telegram file upload failed with status %d: %s\n", resp.StatusCode, string(bodyBytes))
		}
		return fmt.Errorf("Telegram file upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if !silent {
		fmt.Printf("[INFO] File %s sent to Telegram\n", filepath.Base(filePath))
	}
	return nil
}
