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
			show.Episodes, err = tvmaze.FetchEpisodes(show.ID)

			if err != nil {
				return tvmaze.Show{}, err
			}

			return show, nil
		}
	}

	fmt.Println("Unable to find show.")
	return tvmaze.Show{}, nil
}

func detectUnreleasedEpisodes(show tvmaze.Show) {
	for _, ep := range show.Episodes {
		if !ep.IsReleased() {
			fmt.Printf("Unreleased episode: S%s E%s %s \n", strconv.Itoa(ep.Season), strconv.Itoa(ep.Number), ep.Name)
		}
	}
}

func menu(scanner *bufio.Scanner) {
	fmt.Println("1. Add show")
	fmt.Println("2. View shows")
	fmt.Println("3. Exit")

	var input string
	scanner.Scan()
	input = scanner.Text()

	switch input {
	case "1":
		show, err := searchShow(scanner)

		if err != nil {
			panic(err)
		}

		fmt.Printf("You selected: %+v\n", show)
		detectUnreleasedEpisodes(show)
	case "2":
		fmt.Println("View shows")
	case "3":
		os.Exit(0)
	default:
		fmt.Println("Invalid input, please try again")
	}

	menu(scanner)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	menu(scanner)
}
