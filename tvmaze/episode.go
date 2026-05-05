package tvmaze

import (
	"encoding/json"
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

func (e Episode) IsReleased() bool {
	releaseTime, err := time.Parse(time.RFC3339, e.Airstamp)

	if err != nil {
		panic(err)
	}

	if releaseTime.After(time.Now()) {
		return false
	}

	return true
}

func (e Episode) WasReleasedInLast24Hours() bool {
	releaseTime, err := time.Parse(time.RFC3339, e.Airstamp)
	last24Hours := time.Now().Add(-24 * time.Hour)

	if err != nil {
		panic(err)
	}

	if releaseTime.After(last24Hours) && releaseTime.Before(time.Now()) {
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
