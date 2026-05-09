package telegram

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type Client struct{}

func (c Client) SendMessage(message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		return fmt.Errorf("missing TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID environment variable")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	data := fmt.Sprintf("chat_id=%s&text=%s", chatID, message)

	_, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))

	if err != nil {
		slog.Error("Failed to send Telegram message", slog.String("error", err.Error()))
		return err
	}

	return nil
}
