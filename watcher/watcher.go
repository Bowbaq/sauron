package watcher

import (
	"time"

	"github.com/Bowbaq/sauron/model"
)

// Watcher is the interface required by Sauron to detect changes.
type Watcher interface {
	CheckForUpdate(opts model.WatchOptions, since time.Time) (model.Update, error)
}
