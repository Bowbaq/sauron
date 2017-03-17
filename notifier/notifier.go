package notifier

import "github.com/Bowbaq/sauron/model"

// Notifier is the interface required by Sauron to publish changes
type Notifier interface {
	// Notify publishes changes detected by Sauron
	Notify(opts model.WatchOptions, lastUpdate, update model.Update) error
}
