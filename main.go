package main

import (
	"bufio"
	"log/slog"
	"os"
	"show-notifier/menu"
	"show-notifier/notifier"
	"show-notifier/storage"
	"show-notifier/telegram"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	scanner := bufio.NewScanner(os.Stdin)
	store, err := storage.LoadOrCreateStore()

	if err != nil {
		panic(err)
	}

	slog.Info("Starting scheduler")
	client := telegram.Client{}
	go notifier.StartScheduler(&store, client, 24*time.Hour)

	if os.Getenv("MODE") == "HEADLESS" {
		slog.Info("Running in headless mode, skipping menu")
		telegram.PollUpdates(&store, client, 0)
	} else {
		slog.Info("Running in normal mode")
		menu.Menu(scanner, store, client)
	}
}
