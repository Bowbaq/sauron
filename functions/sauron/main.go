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
