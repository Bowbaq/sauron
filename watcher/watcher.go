package watcher

import (
	"time"

	"github.com/Bowbaq/sauron/model"
)

// Watcher is the interface required by Sauron to detect changes.
type Watcher interface {
	CheckForUpdate(opts model.WatchOptions, since time.Time) (model.Update, error)
}

// Options contains the parameters necessary to configure a concrete Watcher implementation
type Options struct {
}

// New returns a new concrete Watcher implementation. For now, only Github is supported, it requires no
// special configuration and is the default
func New(opts Options) Watcher {
	return NewGithub()
}
