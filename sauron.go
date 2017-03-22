package sauron

import (
	"github.com/Bowbaq/belt"

	"github.com/Bowbaq/sauron/model"
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

// Options contains the parameters necessary to configure Sauron
type Options struct {
	WatcherOptions  watcher.Options  `group:"watcher" namespace:"watcher"`
	NotifierOptions notifier.Options `group:"notifier" namespace:"notifier"`
	StoreOptions    store.Options    `group:"store" namespace:"store"`
}

// New creates a new instance of Sauron. If no options are provided, the sate is stored in a file (.sauron)
// and notifications are written to standard error
func New(opts Options) *Sauron {
	return &Sauron{
		watcher: watcher.New(opts.WatcherOptions),

		notifier: notifier.New(opts.NotifierOptions),

		db: store.New(opts.StoreOptions),
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

// Watch checks for updates in the target repository
func (s *Sauron) Watch(opts model.WatchOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	lastUpdate, err := s.db.GetLastUpdate(store.Key(opts))
	if err != nil {
		return err
	}

	update, err := s.watcher.CheckForUpdate(opts, lastUpdate.Timestamp)
	if err != nil {
		return err
	}

	if update.IsZero() || update.IsNotAfter(lastUpdate) {
		belt.Debugf(
			"sauron: [%s b: %s, p: %s] no updates since the last run",
			opts.Repository, opts.Branch, opts.Path,
		)
		return s.db.SetLastChecked(store.Key(opts))
	}

	belt.Debugf(
		"sauron: [%s b: %s, p: %s] updated at %v (%6s)",
		opts.Repository, opts.Branch, opts.Path, update.Timestamp, update.SHA,
	)
	err = s.db.RecordUpdate(store.Key(opts), update)
	if err != nil {
		return err
	}

	// Only notify if there was a change, not on the first run
	if !lastUpdate.IsZero() {
		return s.notifier.Notify(opts, lastUpdate, update)
	}

	return nil
}
