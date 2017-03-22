package notifier

import (
	"fmt"
	"os"

	"github.com/Bowbaq/sauron/model"
)

type fileNotifier struct {
	path string
}

// NewFile creates a basic Notifier that prints to standard error
func NewFile(path string) Notifier {
	return &fileNotifier{
		path: path,
	}
}

func (fn *fileNotifier) Notify(opts model.WatchOptions, lastUpdate, update model.Update) error {
	var (
		out = os.Stderr
		err error
	)
	if fn.path != "" {
		out, err = os.OpenFile(fn.path, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer out.Close()
	}

	fmt.Fprintf(
		out,
		"%s was updated at %v from %6s to %6s\n",
		opts.Repository, update.Timestamp, lastUpdate.SHA, update.SHA,
	)

	return nil
}
