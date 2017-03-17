package store

import (
	"time"

	"github.com/Bowbaq/sauron/model"
)

type jsonStore struct {
	read  func() (State, error)
	write func(State) error
}

// GetLastUpdate returns the last time a repository was updated.
func (js *jsonStore) GetLastUpdate(key WatchKey) (model.Update, error) {
	state, err := js.read()
	if err != nil {
		return model.Update{}, err
	}

	return state[key].Update, nil
}

// RecordUpdate records the last update for a specific repository.
func (js *jsonStore) RecordUpdate(key WatchKey, update model.Update) error {
	state, err := js.read()
	if err != nil {
		return err
	}

	repoState := state[key]
	repoState.Update = update
	repoState.LastChecked = time.Now().UTC()
	state[key] = repoState

	return js.write(state)
}

// SetLastChecked records the last check time for a specific repository.
func (js *jsonStore) SetLastChecked(key WatchKey) error {
	state, err := js.read()
	if err != nil {
		return err
	}

	repoState := state[key]
	repoState.LastChecked = time.Now().UTC()
	state[key] = repoState

	return js.write(state)
}
