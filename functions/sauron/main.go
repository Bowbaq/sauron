package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Bowbaq/sauron"
	"github.com/Bowbaq/sauron/notifier"
	"github.com/Bowbaq/sauron/store"
	"github.com/apex/go-apex"
)

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		s := sauron.New()
		s.SetStore(store.NewS3(os.Getenv("S3_BUCKET"), os.Getenv("S3_KEY")))
		s.SetNotifier(notifier.NewSNS(os.Getenv("SNS_TOPIC_ARN")))

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
