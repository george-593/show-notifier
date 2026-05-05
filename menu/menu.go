package menu

import (
	"bufio"
	"fmt"
	"os"
	"show-notifier/storage"
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

	store.AddShow(show)

	err = storage.Save(*store)

	if err != nil {
		panic(err)
	}
}

func loadShow(store *storage.Store) {
	for _, show := range store.Shows {
		fmt.Printf("Show: %s\n", show.Name)
	}

	if len(store.Shows) == 0 {
		fmt.Println("No shows added yet.")
	}
}

func removeShow(scanner *bufio.Scanner, store *storage.Store) {
	for i, show := range store.Shows {
		fmt.Printf("%d. %s\n", i+1, show.Name)
	}

	if len(store.Shows) == 0 {
		fmt.Println("No shows added yet.")
		return
	}

	fmt.Print("Enter the number of the show you want to remove: ")
	var input string
	scanner.Scan()
	input = scanner.Text()

	index, err := strconv.Atoi(input)

	if err != nil || index < 1 || index > len(store.Shows) {
		fmt.Println("Invalid input, please try again")
		return
	}

	store.Shows = append(store.Shows[:index-1], store.Shows[index:]...)

	err = storage.Save(*store)

	if err != nil {
		panic(err)
	}

	fmt.Println("Show removed successfully.")
}

func Menu(scanner *bufio.Scanner, store storage.Store) {
	for {
		fmt.Println("1. Add show")
		fmt.Println("2. View shows")
		fmt.Println("3. Remove show")
		fmt.Println("4. Exit")

		var input string
		scanner.Scan()
		input = scanner.Text()

		switch input {
		case "1":
			addShow(scanner, &store)
		case "2":
			loadShow(&store)
		case "3":
			removeShow(scanner, &store)
		case "4":
			os.Exit(0)
		default:
			fmt.Println("Invalid input, please try again")
		}

	}
}
