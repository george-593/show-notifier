package tvmaze

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Episode struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	Season   int    `json:"season"`
	Number   int    `json:"number"`
	Type     string `json:"type"`
	Airdate  string `json:"airdate"`
	Airtime  string `json:"airtime"`
	Airstamp string `json:"airstamp"`
	Runtime  int    `json:"runtime"`
	Summary  string `json:"summary"`
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
