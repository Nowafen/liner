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
	"sync"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorWhite  = "\033[37;1m"
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
			fmt.Printf("%s[WARNING]%s Failed to create Telegram message request: %v%s\n",
				ColorYellow, ColorWhite, err, ColorReset)
		}
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{} // No timeout
	resp, err := client.Do(req)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to send Telegram message: %v%s\n",
				ColorYellow, ColorWhite, err, ColorReset)
		}
		return fmt.Errorf("failed to send telegram message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if !silent {
			fmt.Printf("%s[WARNING]%s Telegram message failed with status %d: %s%s\n",
				ColorYellow, ColorWhite, resp.StatusCode, string(bodyBytes), ColorReset)
		}
		return fmt.Errorf("telegram message failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if !silent {
		fmt.Printf("%s[INFO]%s Message sent to Telegram: %s%s\n",
			ColorGreen, ColorWhite, message, ColorReset)
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
	_ = writer.WriteField("chat_id", chatID)
	if err := writer.Close(); err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to close form data for %s: %v%s\n",
				ColorYellow, ColorWhite, filePath, err, ColorReset)
		}
		return fmt.Errorf("failed to close form data: %v", err)
	}

	// Create and send request
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", token)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to create Telegram request for %s: %v%s\n",
				ColorYellow, ColorWhite, filePath, err, ColorReset)
		}
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{} // No timeout
	resp, err := client.Do(req)
	if err != nil {
		if !silent {
			fmt.Printf("%s[WARNING]%s Failed to send file %s to Telegram: %v%s\n",
				ColorYellow, ColorWhite, filePath, err, ColorReset)
		}
		return fmt.Errorf("failed to send file to telegram: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if !silent {
			fmt.Printf("%s[WARNING]%s Telegram file upload failed with status %d: %s%s\n",
				ColorYellow, ColorWhite, resp.StatusCode, string(bodyBytes), ColorReset)
		}
		return fmt.Errorf("telegram file upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if !silent {
		fmt.Printf("%s[INFO]%s File %s sent to Telegram%s\n",
			ColorGreen, ColorWhite, filepath.Base(filePath), ColorReset)
	}
	return nil
}

// SendFilesConcurrently sends multiple files to Telegram concurrently
func SendFilesConcurrently(token, chatID string, filePaths []string, silent bool) []error {
	var wg sync.WaitGroup
	errors := make([]error, 0, len(filePaths))
	errorChan := make(chan error, len(filePaths))

	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			if err := SendToTelegram(token, chatID, filePath, silent); err != nil {
				errorChan <- fmt.Errorf("failed to send file %s: %v", filePath, err)
			}
		}(filePath)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChan)

	// Collect errors
	for err := range errorChan {
		errors = append(errors, err)
	}

	return errors
}
