package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	show, err := fetchShow(84190)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Show: %+v\n", show)
	fmt.Println(show.Name)

	show.Episodes, err = fetchEpisodes(show.ID)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Episodes: %+v\n", show.Episodes)
}
