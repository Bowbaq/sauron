package store

import (
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

type sauronState map[string]repoState

type repoState struct {
	LastUpdated time.Time
	LastCommit  string
	LastChecked time.Time
}

type jsonStore struct {
	read  func() (sauronState, error)
	write func(sauronState) error
}

// GetLastUpdated returns the last time a repository was updated.
func (js *jsonStore) GetLastUpdated(owner, repo string) (time.Time, string, error) {
	state, err := js.read()
	if err != nil {
		return time.Time{}, "", err
	}

	key := fmt.Sprintf("%s/%s", owner, repo)
	repoState := state[key]

	return repoState.LastUpdated, repoState.LastCommit, nil
}

// SetLastUpdated records the last update for a specific repository.
func (js *jsonStore) SetLastUpdated(owner, repo string, commit *github.Commit) error {
	state, err := js.read()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", owner, repo)

	repoState := state[key]
	repoState.LastUpdated = *commit.Author.Date
	repoState.LastCommit = *commit.Tree.SHA
	repoState.LastChecked = time.Now().UTC()
	state[key] = repoState

	return js.write(state)
}

// SetLastChecked records the last check time for a specific repository.
func (js *jsonStore) SetLastChecked(owner, repo string) error {
	state, err := js.read()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", owner, repo)

	repoState := state[key]
	repoState.LastChecked = time.Now().UTC()
	state[key] = repoState

	return js.write(state)
}
