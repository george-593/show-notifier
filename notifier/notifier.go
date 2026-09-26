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

	resp, err := http.Get("https://api.tvmaze.com/updates/shows?since=week")

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

	changed := false

	// The feed covers a week so the same shows reappear for days. Only refetch
	// a show whose timestamp is newer than the one already stored.
	for updatedShowID, updatedAt := range updates {
		for i, show := range store.Shows {
			if show.ID == updatedShowID && updatedAt > show.Updated {
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
				changed = true
				slog.Info("Show info updated successfully", slog.String("show_name", updatedShow.Name))

			}
		}
	}

	if !changed {
		slog.Info("No shows changed, skipping save")
		return
	}

	err = storage.Save(*store)

	if err != nil {
		slog.Error("Failed to save store after processing updates", slog.String("error", err.Error()))
	} else {
		slog.Info("Store saved successfully after processing updates")
	}
}

func untilNextMidnight(now time.Time, loc *time.Location) time.Duration {
	now = now.In(loc)
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
	return next.Sub(now)
}

// StartScheduler runs two independent schedules on a single goroutine:
// episode detection every detectInterval (purely local, no network), and a
// refetch of show data from TVMaze at 00:00 in loc.
func StartScheduler(store *storage.Store, n Notifier, detectInterval time.Duration, loc *time.Location) {
	FetchUpdates(store)
	DetectNewEpisodes(store, n)

	detect := time.NewTicker(detectInterval)
	defer detect.Stop()

	refetch := time.NewTimer(untilNextMidnight(time.Now(), loc))
	defer refetch.Stop()

	for {
		select {
		case <-detect.C:
			slog.Info("Running scheduled episode detection")
			DetectNewEpisodes(store, n)
		case <-refetch.C:
			slog.Info("Running scheduled show refetch")
			FetchUpdates(store)
			// An episode pulled in by the refetch may already have aired, so
			// detect straight away rather than waiting for the next tick.
			DetectNewEpisodes(store, n)
			refetch.Reset(untilNextMidnight(time.Now(), loc))
		}
	}
}
