package storage

import (
	"encoding/json"
	"os"
	"show-notifier/tvmaze"
	"time"
)

var StorePath string = "store.json"

type Store struct {
	Updated     time.Time
	Shows       []tvmaze.Show
	NotifiedIDs []int
}

func (s *Store) AddShow(show tvmaze.Show) {
	s.Shows = append(s.Shows, show)
	s.Updated = time.Now()
}

func (s *Store) ContainsShow(show tvmaze.Show) bool {
	for _, s := range s.Shows {
		if s.ID == show.ID {
			return true
		}
	}

	return false
}

func (s *Store) ContainsNotifiedID(episodeID int) bool {
	for _, id := range s.NotifiedIDs {
		if id == episodeID {
			return true
		}
	}

	return false
}

func (s *Store) MarkNotified(episodeID int) {
	s.NotifiedIDs = append(s.NotifiedIDs, episodeID)
	s.Updated = time.Now()
}

func createStore(shows []tvmaze.Show) Store {
	return Store{
		Updated:     time.Now(),
		Shows:       shows,
		NotifiedIDs: []int{},
	}
}

func Save(store Store) error {
	data, err := json.Marshal(store)

	if err != nil {
		return err
	}

	return os.WriteFile(StorePath, data, 0644)
}

func LoadOrCreateStore() (Store, error) {
	data, err := os.ReadFile(StorePath)

	if os.IsNotExist(err) {
		return createStore([]tvmaze.Show{}), nil
	}

	if err != nil {
		return Store{}, err
	}

	var store Store
	err = json.Unmarshal(data, &store)

	if err != nil {
		return Store{}, err
	}

	return store, nil
}
