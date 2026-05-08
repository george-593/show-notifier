package notifier

import (
	"fmt"
	"show-notifier/storage"
	"strconv"
	"time"
)

type Notifier interface {
	SendMessage(message string) error
}

func DetectNewEpisodes(store *storage.Store, n Notifier) error {
	for _, show := range store.Shows {
		for _, ep := range show.Episodes {
			if ep.WasReleasedInLast24Hours() && !store.ContainsNotifiedID(ep.ID) {
				message := fmt.Sprintf("New episode released: %s S%s E%s %s", show.Name, strconv.Itoa(ep.Season), strconv.Itoa(ep.Number), ep.Name)
				err := n.SendMessage(message)

				if err != nil {
					return fmt.Errorf("failed to send notification for new episode: %v", err)
				} else {
					store.MarkNotified(ep.ID)
					err = storage.Save(*store)

					if err != nil {
						return fmt.Errorf("failed to save store after marking episode as notified: %v", err)
					}
				}
			}
		}
	}
	return nil
}

func StartScheduler(store *storage.Store, n Notifier, interval time.Duration) {
	err := DetectNewEpisodes(store, n)

	if err != nil {
		fmt.Printf("Error during initial episode detection: %v\n", err)
	}

	ticker := time.NewTicker(interval)
	for range ticker.C {
		fmt.Println("Running scheduled episode detection") // Change to log when setup
		err := DetectNewEpisodes(store, n)
		if err != nil {
			fmt.Printf("Error during scheduled episode detection: %v\n", err)
		}
	}
}
