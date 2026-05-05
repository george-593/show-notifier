package main

import (
	"bufio"
	"fmt"
	"os"
	"show-notifier/tvmaze"
	"strconv"
	"strings"
)

func searchShow(scanner *bufio.Scanner) (tvmaze.Show, error) {

	var search string
	fmt.Print("Enter the show you are searching for: ")
	scanner.Scan()
	search = scanner.Text()

	if search == "" {
		fmt.Println("Invalid input, please try again")
		return searchShow(scanner)
	}

	res, err := tvmaze.SearchShows(search)

	if err != nil {
		return tvmaze.Show{}, err
	}

	fmt.Println("Got " + strconv.Itoa(len(res)) + " results.")
	for i := 0; i < len(res); i++ {
		var show tvmaze.Show = res[i].Show
		date := "(" + strings.Split(show.Premiered, "-")[0] + ")"

		fmt.Printf("Are you searching for %s %s? (y/n) ", show.Name, date)

		var answer string
		scanner.Scan()
		answer = scanner.Text()

		if answer == "y" {
			return show, nil
		}
	}

	fmt.Println("Unable to find show.")
	return tvmaze.Show{}, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	search, err := searchShow(scanner)

	if err != nil {
		panic(err)
	}

	search.Episodes, err = tvmaze.FetchEpisodes(search.ID)

	if err != nil {
		panic(err)
	}

	fmt.Printf("You selected: %+v\n", search)
}
