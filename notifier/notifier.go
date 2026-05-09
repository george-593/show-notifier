package notifier

import (
	"fmt"
	"log/slog"
	"show-notifier/storage"
	"strconv"
	"time"
)

type Notifier interface {
	SendMessage(message string) error
}

func DetectNewEpisodes(store *storage.Store, n Notifier) {
	slog.Info("Detecting new episodes")
	for _, show := range store.Shows {
		for _, ep := range show.Episodes {
			if ep.WasReleasedInLast24Hours() && !store.ContainsNotifiedID(ep.ID) {
				message := fmt.Sprintf("New episode released: %s S%s E%s %s", show.Name, strconv.Itoa(ep.Season), strconv.Itoa(ep.Number), ep.Name)
				slog.Info("Sending notification for new episode", slog.String("message", message))
				err := n.SendMessage(message)

				if err != nil {
					slog.Error("Failed to send notification for new episode", slog.String("error", err.Error()))
				} else {
					store.MarkNotified(ep.ID)
					slog.Info("Marking episode as notified", slog.Int("episode_id", ep.ID))
					err = storage.Save(*store)

					if err != nil {
						slog.Error("Failed to save store after marking episode as notified", slog.String("error", err.Error()))
					}
				}
			}
		}
	}
}

func StartScheduler(store *storage.Store, n Notifier, interval time.Duration) {
	DetectNewEpisodes(store, n)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		slog.Info("Running scheduled episode detection")
		DetectNewEpisodes(store, n)
	}
}
