package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Show struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Genres    []string `json:"genres"`
	Premiered string   `json:"premiered"`
	Ended     string   `json:"ended"`
	Externals struct {
		Tvrage  string `json:"tvrage"`
		Thetvdb int    `json:"thetvdb"`
		Imdb    string `json:"imdb"`
	} `json:"externals"`
	Updated int `json:"updated"`
	Links   struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Previousepisode struct {
			Href string `json:"href"`
		} `json:"previousepisode"`
	} `json:"_links"`
	Episodes []Episode `json:"-"`
}

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

func fetchShow(id int) (Show, error) {
	resp, err := http.Get("https://api.tvmaze.com/shows/" + strconv.Itoa(id))

	if err != nil {
		return Show{}, err
	}

	defer resp.Body.Close()

	var show Show
	err = json.NewDecoder(resp.Body).Decode(&show)

	if err != nil {
		return Show{}, err
	}

	return show, nil
}

func fetchEpisodes(showID int) ([]Episode, error) {
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
