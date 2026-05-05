package tvmaze

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type SearchResult struct {
	Score float64 `json:"score"`
	Show  Show    `json:"show"`
}

func SearchShows(search string) ([]SearchResult, error) {
	resp, err := http.Get("https://api.tvmaze.com/search/shows?q=" + url.QueryEscape(search))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []SearchResult
	err = json.NewDecoder(resp.Body).Decode(&results)

	if err != nil {
		return nil, err
	}

	return results, nil
}
