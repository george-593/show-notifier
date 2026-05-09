package tvmaze

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Episode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Season   int    `json:"season"`
	Number   int    `json:"number"`
	Airdate  string `json:"airdate"`
	Airtime  string `json:"airtime"`
	Airstamp string `json:"airstamp"`
}

func (e Episode) WasReleasedInLast24Hours() bool {
	releaseTime, err := time.Parse(time.RFC3339, e.Airstamp)
	last24Hours := time.Now().Add(-24 * time.Hour)

	if err != nil {
		slog.Error("Failed to parse episode release time", slog.String("error", err.Error()), slog.Int("episode_id", e.ID))
		return false
	}

	if releaseTime.After(last24Hours) && releaseTime.Before(time.Now()) {
		return true
	}
	return false
}

func (e Episode) WillReleaseInNextWeek() bool {
	releaseTime, err := time.Parse(time.RFC3339, e.Airstamp)
	nextWeek := time.Now().Add(7 * 24 * time.Hour)

	if err != nil {
		slog.Error("Failed to parse episode release time", slog.String("error", err.Error()), slog.Int("episode_id", e.ID))
		return false
	}

	if releaseTime.After(time.Now()) && releaseTime.Before(nextWeek) {
		return true
	}
	return false
}

func FetchEpisodes(showID int) ([]Episode, error) {
	resp, err := http.Get("https://api.tvmaze.com/shows/" + strconv.Itoa(showID) + "/episodes")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var episodes []Episode
	err = json.NewDecoder(resp.Body).Decode(&episodes)

	if err != nil {
		return nil, err
	}

	return episodes, nil
}
