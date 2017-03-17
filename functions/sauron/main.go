package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Bowbaq/belt"
	"github.com/Bowbaq/sauron"
	"github.com/Bowbaq/sauron/notifier"
	"github.com/Bowbaq/sauron/store"
	"github.com/apex/go-apex"
)

var (
	ErrBucketRequired   = errors.New("sauron-lambda: Bucket name cannot be empty")
	ErrKeyRequired      = errors.New("sauron-lambda: Key cannot be empty")
	ErrSNSTopicRequired = errors.New("sauron-lambda: SNS topic cannot be empty")
)

func init() {
	if os.Getenv("DEBUG") != "" {
		belt.Verbose = true
	}
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		bucket, key := os.Getenv("S3_BUCKET"), os.Getenv("S3_KEY")
		if bucket == "" {
			return nil, ErrBucketRequired
		}
		if key == "" {
			return nil, ErrKeyRequired
		}

		snsTopic := os.Getenv("SNS_TOPIC_ARN")
		if snsTopic == "" {
			return nil, ErrSNSTopicRequired
		}

		s := sauron.New()
		s.SetStore(store.NewS3(bucket, key))
		s.SetNotifier(notifier.NewSNS(snsTopic))

		var opts sauron.WatchOptions
		if err := json.Unmarshal(event, &opts); err != nil {
			return nil, err
		}

		if err := s.Watch(&opts); err != nil {
			return nil, fmt.Errorf("sauron-lambda: Error retrieving latest update: %v", err)
		}

		return nil, nil
	})
}
