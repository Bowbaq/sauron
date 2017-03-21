package notifier

import (
	"fmt"
	"os"

	"github.com/Bowbaq/sauron/model"
)

type stderrNotifier struct{}

// NewStderr creates a basic Notifier that prints to standard error
func NewStderr() Notifier {
	return &stderrNotifier{}
}

func (sn *stderrNotifier) Notify(opts model.WatchOptions, lastUpdate, update model.Update) error {
	fmt.Fprintf(
		os.Stderr,
		"%s was updated at %v from %6s to %6s\n",
		opts.Repository, update.Timestamp, lastUpdate.SHA, update.SHA,
	)

	return nil
}
