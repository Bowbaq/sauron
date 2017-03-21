package watcher

import (
	"context"
	"time"

	"github.com/google/go-github/github"

	"github.com/Bowbaq/sauron/model"
)

type githubWatcher struct {
	ghClient *github.Client
}

// NewGithub instantiates a new watcher for public GitHub repositories.
func NewGithub() Watcher {
	return &githubWatcher{
		ghClient: github.NewClient(nil),
	}
}

// LastCommit retrieves the last commit for the specified repository, optionally restricting the results by
// branch, path and/or updates after a given time.
func (gw *githubWatcher) CheckForUpdate(opts model.WatchOptions, since time.Time) (model.Update, error) {
	listOpts := &github.CommitsListOptions{
		Since: since,

		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	}
	if opts.Branch != "" {
		listOpts.SHA = opts.Branch
	}
	if opts.Path != "" {
		listOpts.Path = opts.Path
	}

	commits, _, err := gw.ghClient.Repositories.ListCommits(
		context.Background(), opts.Repository.Owner, opts.Repository.Name, listOpts,
	)
	if err != nil {
		return model.Update{}, err
	}

	if len(commits) == 0 {
		return model.Update{}, nil
	}

	commit := commits[0].Commit

	return model.Update{
		Timestamp: *commit.Author.Date,
		SHA:       *commit.Tree.SHA,
	}, nil
}
