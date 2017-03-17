package store

import (
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	// Import the PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/Bowbaq/sauron/model"
)

const schema = `
CREATE TABLE IF NOT EXISTS state (
  owner        varchar,
  name         varchar,
  branch       varchar DEFAULT '',
  path         varchar DEFAULT '',
  timestamp    timestamp,
  sha          varchar(40),
  last_checked timestamp,
  CONSTRAINT key PRIMARY KEY(owner, name, branch, path)
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

// GetLastUpdate returns the last time a repository was updated.
func (ps *postgresStore) GetLastUpdate(key WatchKey) (model.Update, error) {
	rows, err := ps.db.NamedQuery(`
    SELECT timestamp, sha FROM state
    WHERE owner = :repository.owner AND name = :repository.name AND branch = :branch AND path = :path`,
		key,
	)
	if err != nil {
		return model.Update{}, err
	}

	var u model.Update
	rows.Next()
	err = rows.StructScan(&u)
	if err != nil {
		if err.Error() == "sql: Rows are closed" {
			return u, nil
		}
		return model.Update{}, err
	}

	return u, nil
}

// RecordUpdate records the last update for a specific repository.
func (ps *postgresStore) RecordUpdate(key WatchKey, update model.Update) error {
	_, err := ps.db.NamedExec(`
    INSERT INTO state (
      owner, name, branch, path, timestamp, sha, last_checked
    )
    VALUES (
      :repository.owner, :repository.name, :branch, :path, :timestamp, :sha, :last_checked
    )
    ON CONFLICT (owner, name, branch, path) DO UPDATE SET
      timestamp = :timestamp, sha = :sha, last_checked = :last_checked`,
		struct {
			WatchKey
			RepoState
		}{
			WatchKey: key,
			RepoState: RepoState{
				Update:      update,
				LastChecked: time.Now().UTC(),
			},
		},
	)

	return err
}

// SetLastChecked records the last check time for a specific repository.
func (ps *postgresStore) SetLastChecked(key WatchKey) error {
	_, err := ps.db.NamedExec(
		`UPDATE state SET last_checked = :last_checked
     WHERE owner = :repository.owner AND name = :repository.name AND branch = :branch AND path = :path`,
		struct {
			WatchKey
			RepoState
		}{
			WatchKey: key,
			RepoState: RepoState{
				LastChecked: time.Now().UTC(),
			},
		},
	)

	return err
}
