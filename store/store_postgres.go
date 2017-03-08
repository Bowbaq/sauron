package store

import (
	"database/sql"
	"os"
	"time"

	// Import the PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/google/go-github/github"
	"github.com/jmoiron/sqlx"
)

const schema = `
CREATE TABLE IF NOT EXISTS state (
    owner               varchar(128),
    repo                varchar(128),
    last_updated        timestamp,
    last_commit         varchar(40),
    last_checked        timestamp,
    CONSTRAINT owner_repo PRIMARY KEY(owner, repo)
);
`

const dropSchema = `
DROP TABLE IF EXISTS state;
`

type postgresStore struct {
	db *sqlx.DB
}

// NewPostgres instantiates a new concrete Store. NewPostgres panics if the dataSource is invalid
// or if the schema cannot be initialized.
func NewPostgres(dataSource string) Store {
	ps := &postgresStore{
		db: sqlx.MustConnect("postgres", dataSource),
	}
	if os.Getenv("DROP_TABLES") != "" {
		ps.db.MustExec(dropSchema)
	}
	ps.db.MustExec(schema)

	return ps
}

// GetLastUpdated returns the last time a repository was updated.
func (ps *postgresStore) GetLastUpdated(owner, repo string) (time.Time, string, error) {
	var result struct {
		LastUpdated time.Time `db:"last_updated"`
		LastCommit  string    `db:"last_commit"`
	}
	err := ps.db.Get(&result, "SELECT last_updated, last_commit FROM state WHERE owner = $1 AND repo = $2", owner, repo)
	if err == sql.ErrNoRows {
		err = nil
	}

	return result.LastUpdated, result.LastCommit, err
}

// SetLastUpdated records the last update for a specific repository.
func (ps *postgresStore) SetLastUpdated(owner, repo string, commit *github.Commit) error {
	_, err := ps.db.NamedExec(`
    INSERT INTO state (
      owner, repo, last_updated, last_commit, last_checked
    )
    VALUES (
      :owner, :repo, :last_updated, :last_commit, :last_checked
    )
    ON CONFLICT (owner, repo) DO UPDATE SET
      last_updated = :last_updated, last_commit = :last_commit, last_checked = :last_checked`,
		map[string]interface{}{
			"owner":        owner,
			"repo":         repo,
			"last_updated": commit.Author.Date,
			"last_commit":  commit.Tree.SHA,
			"last_checked": time.Now().UTC(),
		},
	)

	return err
}

// SetLastChecked records the last check time for a specific repository.
func (ps *postgresStore) SetLastChecked(owner, repo string) error {
	_, err := ps.db.NamedExec(
		`UPDATE state SET last_checked = :last_checked WHERE owner = :owner and repo = :repo`,
		map[string]interface{}{
			"owner":        owner,
			"repo":         repo,
			"last_checked": time.Now().UTC(),
		},
	)

	return err
}
