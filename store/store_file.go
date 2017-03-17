package store

import (
	"encoding/json"
	"io"
	"os"
)

const storePath = ".sauron"

type fileStore struct {
	*jsonStore
}

// NewFile instantiates a new concrete Store backed by a local file.
func NewFile() Store {
	return &fileStore{&jsonStore{
		read: func() (State, error) {
			f, err := os.OpenFile(storePath, os.O_RDONLY|os.O_CREATE, 0644)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			var s State
			err = json.NewDecoder(f).Decode(&s)
			if err != nil {
				if err == io.EOF {
					return make(State), nil
				}
				return nil, err
			}

			return s, nil
		},

		write: func(s State) error {
			f, err := os.OpenFile(storePath, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			return json.NewEncoder(f).Encode(s)
		},
	}}
}
