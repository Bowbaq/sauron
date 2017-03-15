package notifier

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/go-github/github"
)

type snsNotifier struct {
	targetARN string

	client *sns.SNS
}

// NewSNS creates a Notifier that publishes events to an SNS topic
func NewSNS(targetARN string) Notifier {
	return &snsNotifier{
		targetARN: targetARN,

		client: sns.New(session.Must(session.NewSession())),
	}
}

func (sn *snsNotifier) Notify(owner, repo, lastSHA string, newCommit *github.Commit) error {
	message := fmt.Sprintf(
		"%s/%s was updated at %v from %6s to %6s\n",
		owner, repo, *newCommit.Author.Date, lastSHA, *newCommit.Tree.SHA,
	)
	subject := fmt.Sprintf("sauron: change detected %s/%s", owner, repo)

	_, err := sn.client.Publish(&sns.PublishInput{
		Message:   &message,
		Subject:   &subject,
		TargetArn: &sn.targetARN,
	})

	return err
}
