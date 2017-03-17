package store

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Bowbaq/sauron/model"
)

// Store is the storage layer interface for Sauron
type Store interface {
	GetLastUpdate(key WatchKey) (model.Update, error)
	RecordUpdate(key WatchKey, update model.Update) error
	SetLastChecked(key WatchKey) error
}

// WatchKey is used to retrieve a RepoState from the global state
type WatchKey struct {
	Repository model.Repository

	Branch string

	Path string
}

// MarshalText implements the encoding.TextMarshaller interface
func (k WatchKey) MarshalText() ([]byte, error) {
	var b bytes.Buffer

	fmt.Fprintf(&b, "%s|%s|%s|%s", k.Repository.Owner, k.Repository.Name, k.Branch, k.Path)

	return b.Bytes(), nil
}

// UnmarshalText implements the encoding.TextUnmarshaller interface
func (k *WatchKey) UnmarshalText(data []byte) error {
	parts := bytes.Split(data, []byte("|"))
	if len(parts) != 4 {
		return fmt.Errorf("Invalid key %s, expected 4 parts", string(data))
	}

	k.Repository.Owner = string(parts[0])
	k.Repository.Name = string(parts[1])
	k.Branch = string(parts[2])
	k.Path = string(parts[3])

	return nil
}

// Key creates a new WatchKey from a WatchOptions value
func Key(opts model.WatchOptions) WatchKey {
	return WatchKey(opts)
}

// RepoState records the last known state of a given repository
type RepoState struct {
	model.Update

	LastChecked time.Time `db:"last_checked"`
}

// State represents the state tracked by Sauron
type State map[WatchKey]RepoState
