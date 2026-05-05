package telegram

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func SendMessage(message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		return fmt.Errorf("missing TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID environment variable")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	data := fmt.Sprintf("chat_id=%s&text=%s", chatID, message)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return err
	}

	return nil
}
