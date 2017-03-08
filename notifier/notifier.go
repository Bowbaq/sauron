package notifier

import "github.com/google/go-github/github"

type Notifier interface {
	Notify(owner, repo, lastSHA string, newCommit *github.Commit) error
}
