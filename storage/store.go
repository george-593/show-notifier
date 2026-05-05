package storage

import (
	"encoding/json"
	"os"
	"show-notifier/tvmaze"
	"time"
)

type Store struct {
	Updated time.Time
	Shows   []tvmaze.Show
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

func createStore(shows []tvmaze.Show) Store {
	return Store{
		Updated: time.Now(),
		Shows:   shows,
	}
}

func Save(store Store, path string) error {
	data, err := json.Marshal(store)

	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadOrCreateStore(path string) (Store, error) {
	data, err := os.ReadFile(path)

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
