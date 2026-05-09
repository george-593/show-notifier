package telegram

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"show-notifier/notifier"
	"show-notifier/storage"
	"show-notifier/tvmaze"
	"strconv"
	"strings"
	"time"
)

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type User struct {
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type Chat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Type      string `json:"type"`
}

type State struct {
	// Add command state
	AwaitingShowSelection bool
	SearchResults         []tvmaze.SearchResult
}

var state = State{}

func PollUpdates(store *storage.Store, n notifier.Notifier, offset int) {
	for {
		resp, err := http.Get("https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/getUpdates?offset=" + strconv.Itoa(offset))

		if err != nil {
			slog.Error("Failed to get updates", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var updateResp UpdateResponse
		err = json.NewDecoder(resp.Body).Decode(&updateResp)

		if err != nil {
			slog.Error("Failed to decode update response", "error", err)
			resp.Body.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		resp.Body.Close()

		for _, update := range updateResp.Result {
			if update.Message != nil {
				handleMessage(store, n, update.Message)
				offset = update.UpdateID + 1
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func handleMessage(store *storage.Store, n notifier.Notifier, message *Message) {
	slog.Info("Received message", "chat_id", message.Chat.Id, "text", message.Text)

	command, args, _ := strings.Cut(message.Text, " ")

	// Add callback handling
	if state.AwaitingShowSelection {
		handleAddCallback(message.Text, n, store)
		return
	}

	switch command {
	case "/start":
		n.SendMessage("Welcome to the Show Notifier Bot! Use /add to add a show to your watchlist.")
	case "/add":
		handleAdd(n, args)
	default:
		n.SendMessage("Unknown command. Try /add, /list, /remove, /upcoming, /check")
	}
}

func handleAdd(n notifier.Notifier, args string) {
	if args == "" {
		n.SendMessage("Invalid input, please try adding a show to search for")
		slog.Info("Received invalid add command", "args", args)
		return
	}

	res, err := tvmaze.SearchShows(args)

	if err != nil {
		slog.Error("Failed to search for show", slog.String("error", err.Error()))
		n.SendMessage("Failed to search for show, please try again later.")
		return
	}

	message := "Got " + strconv.Itoa(len(res)) + " results.\n"
	for i := 0; i < len(res); i++ {
		var show tvmaze.Show = res[i].Show
		date := "(" + strings.Split(show.Premiered, "-")[0] + ")"

		message += strconv.Itoa(i+1) + ". " + show.Name + " " + date + "\n"
	}

	message += "Reply with the number of the show you want to add."

	n.SendMessage(message)
	slog.Info("Sent search results", "count", len(res))

	state.AwaitingShowSelection = true
	state.SearchResults = res
}

func handleAddCallback(id string, n notifier.Notifier, store *storage.Store) {
	index, err := strconv.Atoi(id)

	if err != nil || index < 1 || index > len(state.SearchResults) {
		n.SendMessage("Invalid selection, please try again.")
		slog.Info("Received invalid show selection", "input", id)
		return
	}

	show := state.SearchResults[index].Show

	if store.ContainsShow(show) {
		n.SendMessage("Show is already in your watchlist.")
		slog.Info("Show already exists in watchlist", "show_id", show.ID)
		return
	}

	if show.Ended != "" {
		n.SendMessage("This show has already ended, adding it anyway.")
		slog.Info("Show has already ended", "show_id", show.ID)
	}

	store.AddShow(show)
	n.SendMessage(fmt.Sprintf("Added %s to your watchlist!", show.Name))
	slog.Info("Added show to watchlist", "show_id", show.ID)

	state.AwaitingShowSelection = false
	state.SearchResults = nil

	saveErr := storage.Save(*store)

	if saveErr != nil {
		slog.Error("Failed to save store after adding show", "error", saveErr)
		n.SendMessage("An error occurred while saving your watchlist. Please try again later.")
	}

}
