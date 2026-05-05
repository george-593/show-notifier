package tvmaze

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Show struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Premiered string `json:"premiered"`
	Ended     string `json:"ended"`
	Externals struct {
		Tvrage  int    `json:"tvrage"`
		Thetvdb int    `json:"thetvdb"`
		Imdb    string `json:"imdb"`
	} `json:"externals"`
	Updated  int       `json:"updated"`
	Episodes []Episode `json:"episodes,omitempty"`
}

func FetchShow(id int) (Show, error) {
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
