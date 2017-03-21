package model

import (
	"fmt"
	"time"

	"github.com/Bowbaq/sauron/errorx"
)

// A Repository is the object being watched by Sauron
type Repository struct {
	Owner string `db:"owner" long:"owner" env:"GITHUB_OWNER" description:"owner of the repository" required:"yes"`

	Name string `db:"name" long:"repository" env:"GITHUB_REPOSITORY" description:"name of the repository" required:"yes"`
}

func (r Repository) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

// WatchOptions specifies where to look for updates
type WatchOptions struct {
	Repository Repository

	// Restrict to a specific branch
	Branch string `long:"branch" env:"GITHUB_BRANCH" description:"branch to watch in the repository"`

	// Restrict to a specific path
	Path string `long:"path" env:"GITHUB_PATH" description:"path to watch in the repository"`
}

// Validate returns an error when options are invalid, nil otherwise
func (opts WatchOptions) Validate() error {
	if opts.Repository.Owner == "" {
		return errorx.ErrRepositoryOwnerRequired
	}

	if opts.Repository.Name == "" {
		return errorx.ErrRepositoryNameRequired
	}

	return nil
}

// Update represents an update to a repository
type Update struct {
	// The timestamp of the commit or release
	Timestamp time.Time `db:"timestamp"`

	// The hash of the commit or release
	SHA string `db:"sha"`
}

// IsNotAfter returns true if an update happened before or at the same time as an other update
func (u Update) IsNotAfter(v Update) bool {
	return !u.Timestamp.After(v.Timestamp)
}

// IsZero returns true when called on the zero-value for Update
func (u Update) IsZero() bool {
	return u.Timestamp.IsZero() && u.SHA == ""
}
