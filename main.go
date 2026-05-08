package main

import (
	"bufio"
	"fmt"
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

	scanner := bufio.NewScanner(os.Stdin)
	store, err := storage.LoadOrCreateStore()

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting scheduler")
	go notifier.StartScheduler(&store, telegram.Client{}, 6*time.Hour)

	menu.Menu(scanner, store)
}
