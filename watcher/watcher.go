package watcher

import (
	"time"

	"github.com/google/go-github/github"
)

// WatchOptions specifies where to look for updates.
type WatchOptions struct {
	// Owner. Required
	Owner string

	// Repository. Required
	Repository string

	// Restrict to a specific branch
	Branch string

	// Restrict to a specific path
	Path string

	// Restrict to updates after a specific date
	Since time.Time
}

// Watcher is the interface required by Sauron to detect changes
type Watcher interface {
	LastCommit(opts *WatchOptions) (*github.Commit, error)
}
