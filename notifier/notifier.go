package notifier

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"show-notifier/storage"
	"show-notifier/tvmaze"
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

func FetchUpdates(store *storage.Store) {
	slog.Info("Fetching updates for all shows")

	resp, err := http.Get("https://api.tvmaze.com/updates/shows?since=day")

	if err != nil {
		slog.Error("Failed to fetch updates from TVMaze API", slog.String("error", err.Error()))
		return
	}

	defer resp.Body.Close()

	var updates map[int]int
	err = json.NewDecoder(resp.Body).Decode(&updates)

	if err != nil {
		slog.Error("Failed to decode updates from TVMaze API", slog.String("error", err.Error()))
		return
	}

	for updatedShowID := range updates {
		for i, show := range store.Shows {
			if show.ID == updatedShowID {
				slog.Info("Updates returned for show, fetching updated info", slog.String("show_name", show.Name))
				updatedShow, err := tvmaze.FetchShow(show.ID)

				if err != nil {
					slog.Error("Failed to fetch updated show info from TVMaze API", slog.String("error", err.Error()))
					continue
				}

				updatedShow.Episodes, err = tvmaze.FetchEpisodes(show.ID)

				if err != nil {
					slog.Error("Failed to fetch updated episodes for show from TVMaze API", slog.String("error", err.Error()))
					continue
				}

				store.Shows[i] = updatedShow
				slog.Info("Show info updated successfully", slog.String("show_name", updatedShow.Name))

			}
		}
	}

	err = storage.Save(*store)

	if err != nil {
		slog.Error("Failed to save store after processing updates", slog.String("error", err.Error()))
	} else {
		slog.Info("Store saved successfully after processing updates")
	}
}

func StartScheduler(store *storage.Store, n Notifier, interval time.Duration) {
	FetchUpdates(store)
	DetectNewEpisodes(store, n)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		slog.Info("Running scheduled episode detection")
		FetchUpdates(store)
		DetectNewEpisodes(store, n)
	}
}
