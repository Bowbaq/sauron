package notifier

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/Bowbaq/sauron/model"
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

func (sn *snsNotifier) Notify(opts model.WatchOptions, lastUpdate, update model.Update) error {
	message := fmt.Sprintf(
		"%s was updated at %v from %6s to %6s\n",
		opts.Repository, update.Timestamp, lastUpdate.SHA, update.SHA,
	)
	subject := fmt.Sprintf("sauron: change detected %s", opts.Repository)

	_, err := sn.client.Publish(&sns.PublishInput{
		Message:   &message,
		Subject:   &subject,
		TargetArn: &sn.targetARN,
	})

	return err
}
