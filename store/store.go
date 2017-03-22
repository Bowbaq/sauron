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

// Options contains the parameters necessary to configure a concrete Store implementation
type Options struct {
	// S3 backend options
	S3 struct {
		Bucket string `long:"bucket" env:"S3_BUCKET" description:"name of the bucket"`
		Key    string `long:"key" default:"state.json" env:"S3_KEY" description:"path to the key"`
	} `group:"store.s3" namespace:"s3"`

	// Postgres backend options
	Postgres struct {
		Datasource string `long:"datasource" env:"PG_DATASOURCE" description:"postgresql datasource (see database/sql)"`
	} `group:"store.postgres" namespace:"pg"`

	// File backend options. This is the default if nothing else is provided
	File struct {
		Path string `long:"path" default:".sauron" env:"STORE_FILE_PATH" description:"path of the state file"`
	} `group:"store.file" namespace:"file"`
}

// New returns a new concrete Store implementation depending on which options are set. If none, a basic
// file-backed implementation is returned
func New(opts Options) Store {
	switch {
	case opts.S3.Bucket != "":
		return NewS3(opts.S3.Bucket, opts.S3.Key)

	case opts.Postgres.Datasource != "":
		return NewPostgres(opts.Postgres.Datasource)

	default:
		return NewFile(opts.File.Path)
	}
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
