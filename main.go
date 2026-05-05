package main

import (
	"bufio"
	"fmt"
	"os"
	"show-notifier/storage"
	"show-notifier/telegram"
	"show-notifier/tvmaze"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var StorePath string = "store.json"

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

func addShow(scanner *bufio.Scanner, store *storage.Store) {
	show, err := searchShow(scanner)

	if err != nil {
		panic(err)
	}

	fmt.Printf("You selected: %+v\n", show.Name)

	if store.ContainsShow(show) {
		fmt.Println("Show already exists.")
		return
	}

	detectUnreleasedEpisodes(show)

	store.AddShow(show)

	err = storage.Save(*store, StorePath)

	if err != nil {
		panic(err)
	}
}

func loadShow(store *storage.Store) {
	for _, show := range store.Shows {
		fmt.Printf("Show: %s\n", show.Name)
		detectUnreleasedEpisodes(show)
	}

	if len(store.Shows) == 0 {
		fmt.Println("No shows added yet.")
	}
}

func detectUnreleasedEpisodes(show tvmaze.Show) {
	for _, ep := range show.Episodes {
		if !ep.IsReleased() {
			fmt.Printf("Unreleased episode: S%s E%s %s \n", strconv.Itoa(ep.Season), strconv.Itoa(ep.Number), ep.Name)
		}
	}
}

func detectNewEpisodes(store storage.Store) {
	for _, show := range store.Shows {
		for _, ep := range show.Episodes {
			if ep.WasReleasedInLast24Hours() {
				message := fmt.Sprintf("New episode released: %s S%s E%s %s", show.Name, strconv.Itoa(ep.Season), strconv.Itoa(ep.Number), ep.Name)
				err := telegram.SendMessage(message)

				if err != nil {
					fmt.Printf("Failed to send message for %s S%s E%s: %v\n", show.Name, strconv.Itoa(ep.Season), strconv.Itoa(ep.Number), err)
				} else {
					fmt.Printf("Sent notification for new episode: %s S%s E%s\n", show.Name, strconv.Itoa(ep.Season), strconv.Itoa(ep.Number))
				}
			}
		}
	}
}

func menu(scanner *bufio.Scanner, store storage.Store) {
	for {
		fmt.Println("1. Add show")
		fmt.Println("2. View shows")
		fmt.Println("3. Exit")

		var input string
		scanner.Scan()
		input = scanner.Text()

		switch input {
		case "1":
			addShow(scanner, &store)
		case "2":
			loadShow(&store)
		case "3":
			os.Exit(0)
		default:
			fmt.Println("Invalid input, please try again")
		}

	}
}

func main() {
	godotenv.Load()

	scanner := bufio.NewScanner(os.Stdin)
	store, err := storage.LoadOrCreateStore(StorePath)

	detectNewEpisodes(store)

	if err != nil {
		panic(err)
	}

	menu(scanner, store)
}
