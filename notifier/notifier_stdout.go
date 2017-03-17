package notifier

import (
	"fmt"

	"github.com/Bowbaq/sauron/model"
)

type stdoutNotifier struct{}

// NewStdout creates a basic Notifier that prints to standard output
func NewStdout() Notifier {
	return &stdoutNotifier{}
}

func (sn *stdoutNotifier) Notify(opts model.WatchOptions, lastUpdate, update model.Update) error {
	fmt.Printf(
		"%s was updated at %v from %6s to %6s\n",
		opts.Repository, update.Timestamp, lastUpdate.SHA, update.SHA,
	)

	return nil
}
