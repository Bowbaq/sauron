package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Bowbaq/belt"
	"github.com/apex/go-apex"

	"github.com/Bowbaq/sauron"
	"github.com/Bowbaq/sauron/flagx"
	"github.com/Bowbaq/sauron/model"
)

var opts sauron.Options

func init() {
	flagx.MustParse(&opts)
}

func init() {
	if os.Getenv("DEBUG") != "" {
		belt.Verbose = true
	}
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		s := sauron.New(opts)

		// Terraform only lets us send flat JSON
		var watchOpts struct {
			Owner      string
			Repository string
			Branch     string
			Path       string
		}
		if err := json.Unmarshal(event, &watchOpts); err != nil {
			return nil, err
		}

		err := s.Watch(model.WatchOptions{
			Repository: model.Repository{
				Owner: watchOpts.Owner,
				Name:  watchOpts.Repository,
			},
			Branch: watchOpts.Branch,
			Path:   watchOpts.Path,
		})
		if err != nil {
			return nil, fmt.Errorf("sauron-lambda: Error retrieving latest update: %v", err)
		}

		return nil, nil
	})
}
