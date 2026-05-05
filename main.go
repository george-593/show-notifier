package main

import (
	"fmt"
	"show-notifier/tvmaze"
)

func main() {
	fmt.Println("Hello, World!")

	show, err := tvmaze.FetchShow(84190)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Show: %+v\n", show)
	fmt.Println(show.Name)

	show.Episodes, err = tvmaze.FetchEpisodes(show.ID)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Episodes: %+v\n", show.Episodes)
}
