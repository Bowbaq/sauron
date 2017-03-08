package watcher

import "github.com/google/go-github/github"

type githubWatcher struct {
	ghClient *github.Client
}

func NewGithub() Watcher {
	return &githubWatcher{
		ghClient: github.NewClient(nil),
	}
}

// LastCommit retrieves the last commit for the specified repository, optionally restricting the results by
// branch, path and/or updates after a given time.
func (gw *githubWatcher) LastCommit(opts *WatchOptions) (*github.Commit, error) {
	listOpts := &github.CommitsListOptions{
		Since: opts.Since,

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

	commits, _, err := gw.ghClient.Repositories.ListCommits(opts.Owner, opts.Repo, listOpts)
	if err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, nil
	}

	return commits[0].Commit, nil
}
