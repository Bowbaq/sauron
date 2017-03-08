package notifier

import (
	"fmt"

	"github.com/google/go-github/github"
)

type stdoutNotifier struct{}

func NewStdout() Notifier {
	return &stdoutNotifier{}
}

func (sn *stdoutNotifier) Notify(owner, repo, lastSHA string, newCommit *github.Commit) error {
	fmt.Printf(
		"%s/%s was updated at %v from %6s to %6s\n",
		owner, repo, *newCommit.Author.Date, lastSHA, *newCommit.Tree.SHA,
	)

	return nil
}
