package store

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	// Import the PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/google/go-github/github"
	"github.com/jmoiron/sqlx"
)

const storePath = ".sauron"

type fileStore struct {
	db *sqlx.DB
}

// NewFile instantiates a new concrete Store backed by a local file.
func NewFile() Store {
	return &fileStore{}
}

// GetLastUpdated returns the last time a repository was updated.
func (fs *fileStore) GetLastUpdated(owner, repo string) (time.Time, string, error) {
	state, err := readStateFile()
	if err != nil {
		return time.Time{}, "", err
	}

	key := fmt.Sprintf("%s/%s", owner, repo)
	repoState := state[key]

	return repoState.LastUpdated, repoState.LastCommit, nil
}

// SetLastUpdated records the last update for a specific repository.
func (fs *fileStore) SetLastUpdated(owner, repo string, commit *github.Commit) error {
	state, err := readStateFile()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", owner, repo)

	repoState := state[key]
	repoState.LastUpdated = *commit.Author.Date
	repoState.LastCommit = *commit.Tree.SHA
	repoState.LastChecked = time.Now().UTC()
	state[key] = repoState

	return writeStateFile(state)
}

// SetLastChecked records the last check time for a specific repository.
func (fs *fileStore) SetLastChecked(owner, repo string) error {
	state, err := readStateFile()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", owner, repo)

	repoState := state[key]
	repoState.LastChecked = time.Now().UTC()
	state[key] = repoState

	return writeStateFile(state)
}

type state map[string]repoState

type repoState struct {
	LastUpdated time.Time
	LastCommit  string
	LastChecked time.Time
}

func readStateFile() (state, error) {
	f, err := os.OpenFile(storePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return make(state), err
	}
	defer f.Close()

	var s state
	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return make(state), err
	}

	return s, nil
}

func writeStateFile(s state) error {
	f, err := os.OpenFile(storePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(s)
}
