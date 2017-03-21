package notifier

import "github.com/Bowbaq/sauron/model"

// Notifier is the interface required by Sauron to publish changes
type Notifier interface {
	// Notify publishes changes detected by Sauron
	Notify(opts model.WatchOptions, lastUpdate, update model.Update) error
}

// Options contains the parameters necessary to configure a concrete Notifier implementation
type Options struct {
	// SNS target options
	SNS struct {
		TopicARN string `long:"topic-arn" env:"SNS_TOPIC_ARN" description:"ARN of the SNS topic"`
	} `group:"notifier.sns" namespace:"sns"`
}

// New returns a new concrete Notifier implementation depending on which options are set. If none, a basic
// notifier writing to stderr is returned
func New(opts Options) Notifier {
	switch {
	case opts.SNS.TopicARN != "":
		return NewSNS(opts.SNS.TopicARN)

	default:
		return NewStderr()
	}
}
