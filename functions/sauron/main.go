package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Bowbaq/belt"
	"github.com/apex/go-apex"

	"github.com/Bowbaq/sauron"
	"github.com/Bowbaq/sauron/errorx"
	"github.com/Bowbaq/sauron/model"
	"github.com/Bowbaq/sauron/notifier"
	"github.com/Bowbaq/sauron/store"
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
			return nil, errorx.ErrBucketNameRequired
		}
		if key == "" {
			return nil, errorx.ErrBucketKeyRequired
		}

		snsTopic := os.Getenv("SNS_TOPIC_ARN")
		if snsTopic == "" {
			return nil, errorx.ErrSNSTopicRequired
		}

		s := sauron.New()
		s.SetStore(store.NewS3(bucket, key))
		s.SetNotifier(notifier.NewSNS(snsTopic))

		var opts model.WatchOptions
		if err := json.Unmarshal(event, &opts); err != nil {
			return nil, err
		}

		if err := s.Watch(opts); err != nil {
			return nil, fmt.Errorf("sauron-lambda: Error retrieving latest update: %v", err)
		}

		return nil, nil
	})
}
