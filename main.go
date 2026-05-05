package main

import (
	"bufio"
	"fmt"
	"os"
	"show-notifier/menu"
	"show-notifier/storage"
	"show-notifier/telegram"
	"strconv"

	"github.com/joho/godotenv"
)

func detectNewEpisodes(store *storage.Store) {
	for _, show := range store.Shows {
		for _, ep := range show.Episodes {
			if ep.WasReleasedInLast24Hours() && !store.ContainsNotifiedID(ep.ID) {
				message := fmt.Sprintf("New episode released: %s S%s E%s %s", show.Name, strconv.Itoa(ep.Season), strconv.Itoa(ep.Number), ep.Name)
				err := telegram.SendMessage(message)

				if err != nil {
					fmt.Printf("Failed to send message for %s S%s E%s: %v\n", show.Name, strconv.Itoa(ep.Season), strconv.Itoa(ep.Number), err)
				} else {
					fmt.Printf("Sent notification for new episode: %s S%s E%s\n", show.Name, strconv.Itoa(ep.Season), strconv.Itoa(ep.Number))
					store.MarkNotified(ep.ID)
					err = storage.Save(*store)

					if err != nil {
						fmt.Printf("Failed to save store after marking episode as notified: %v\n", err)
					}
				}
			}
		}
	}
}

func main() {
	godotenv.Load()

	scanner := bufio.NewScanner(os.Stdin)
	store, err := storage.LoadOrCreateStore()

	if err != nil {
		panic(err)
	}

	detectNewEpisodes(&store)

	menu.Menu(scanner, store)
}
