package store

import (
	"time"

	"github.com/google/go-github/github"
)

// Store is the storage layer interface for Sauron
type Store interface {
	GetLastUpdated(owner, repo string) (time.Time, string, error)
	SetLastUpdated(owner, repo string, commit *github.Commit) error
	SetLastChecked(owner, repo string) error
}
