package main

import (
	"bufio"
	"os"
	"show-notifier/menu"
	"show-notifier/notifier"
	"show-notifier/storage"
	"show-notifier/telegram"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	scanner := bufio.NewScanner(os.Stdin)
	store, err := storage.LoadOrCreateStore()

	if err != nil {
		panic(err)
	}

	err = notifier.DetectNewEpisodes(&store, telegram.Client{})

	if err != nil {
		panic(err)
	}

	menu.Menu(scanner, store)
}
