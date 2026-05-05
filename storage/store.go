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

func CreateStore(shows []tvmaze.Show) Store {
	return Store{
		Updated: time.Now(),
		Shows:   shows,
	}
}

func (s *Store) AddShow(show tvmaze.Show) {
	s.Shows = append(s.Shows, show)
	s.Updated = time.Now()
}

func Save(store Store, path string) error {
	data, err := json.Marshal(store)

	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Load(path string) (Store, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return Store{}, err
	}

	var store Store
	err = json.Unmarshal(data, &store)

	if os.IsNotExist(err) {
		return Store{}, nil
	}

	if err != nil {
		return Store{}, err
	}

	return store, nil
}
