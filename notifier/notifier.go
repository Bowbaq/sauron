package notifier

import "github.com/google/go-github/github"

// Notifier is the interface required by Sauron to publish changes
type Notifier interface {
	// Notify publishes changes detected by Sauron
	Notify(owner, repo, lastSHA string, newCommit *github.Commit) error
}
