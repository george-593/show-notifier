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
}

func fetchShow(id int) Show {
	resp, err := http.Get("https://api.tvmaze.com/shows/" + strconv.Itoa(id))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var show Show
	err = json.NewDecoder(resp.Body).Decode(&show)

	if err != nil {
		panic(err)
	}

	return show
}
