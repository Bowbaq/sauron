package sauron

import (
	"time"

	"github.com/google/go-github/github"

	"github.com/Bowbaq/belt"
	"github.com/Bowbaq/sauron/notifier"
	"github.com/Bowbaq/sauron/store"
	"github.com/Bowbaq/sauron/watcher"
)

// Sauron watches for changes in GitHub repositories
type Sauron struct {
	watcher watcher.Watcher

	notifier notifier.Notifier

	db store.Store
}

// New creates a new instance of Sauron with simple defaults. State is store in a file (.sauron) and
// notifications are written to standard output.
func New() *Sauron {
	return &Sauron{
		watcher: watcher.NewGithub(),

		notifier: notifier.NewStdout(),

		db: store.NewFile(),
	}
}

// SetWatcher overrides the watcher
func (s *Sauron) SetWatcher(w watcher.Watcher) {
	s.watcher = w
}

// SetNotifier overrides the notifier
func (s *Sauron) SetNotifier(n notifier.Notifier) {
	s.notifier = n
}

// SetStore sets the store
func (s *Sauron) SetStore(db store.Store) {
	s.db = db
}

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
}

// Watch checks for updates in the target repository
func (s *Sauron) Watch(opts *WatchOptions) error {
	lastUpdated, lastSHA, err := s.db.GetLastUpdated(opts.Owner, opts.Repository)
	if err != nil {
		return err
	}
	if !lastUpdated.IsZero() {
		belt.Debugf(
			"sauron: [%s/%s b: %s, p: %s] last updated at %v",
			opts.Owner, opts.Repository, opts.Branch, opts.Path, lastUpdated,
		)
	}

	newCommit, err := s.watcher.LastCommit(&watcher.WatchOptions{
		Owner:      opts.Owner,
		Repository: opts.Repository,
		Branch:     opts.Branch,
		Path:       opts.Path,
		Since:      lastUpdated,
	})
	if err != nil {
		return err
	}

	if newCommit == nil || isNotAfter(lastUpdated, newCommit) {
		belt.Debugf(
			"sauron: [%s/%s b: %s, p: %s] no updates since the last run",
			opts.Owner, opts.Repository, opts.Branch, opts.Path,
		)
		return s.db.SetLastChecked(opts.Owner, opts.Repository)
	}

	belt.Debugf(
		"sauron: [%s/%s b: %s, p: %s] updated at %v (%6s)",
		opts.Owner, opts.Repository, opts.Branch, opts.Path, *newCommit.Author.Date, *newCommit.Tree.SHA,
	)
	err = s.db.SetLastUpdated(opts.Owner, opts.Repository, newCommit)
	if err != nil {
		return err
	}

	return s.notifier.Notify(opts.Owner, opts.Repository, lastSHA, newCommit)
}

func isNotAfter(lastUpdated time.Time, commit *github.Commit) bool {
	return !commit.Author.Date.After(lastUpdated)
}
