package store

import (
	"encoding/json"
	"os"
)

const storePath = ".sauron"

type fileStore struct {
	*jsonStore
}

// NewFile instantiates a new concrete Store backed by a local file.
func NewFile() Store {
	return &fileStore{&jsonStore{
		read: func() (sauronState, error) {
			f, err := os.OpenFile(storePath, os.O_RDONLY|os.O_CREATE, 0644)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			var s sauronState
			err = json.NewDecoder(f).Decode(&s)
			if err != nil {
				return nil, err
			}

			return s, nil
		},

		write: func(s sauronState) error {
			f, err := os.OpenFile(storePath, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			return json.NewEncoder(f).Encode(s)
		},
	}}
}
