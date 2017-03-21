package main

import (
	"log"
	"os"

	"github.com/Bowbaq/belt"
	"github.com/Bowbaq/sauron"
	"github.com/Bowbaq/sauron/flagx"
	"github.com/Bowbaq/sauron/model"
)

var (
	// Version of the CLI, filled in at compile time
	Version string

	opts struct {
		sauron.Options
		model.WatchOptions `group:"github" namespace:"github"`
	}
)

func init() {
	flagx.MustParse(&opts)

	if os.Getenv("DEBUG") != "" {
		belt.Verbose = true
	}
}

func main() {
	s := sauron.New(opts.Options)

	if err := s.Watch(opts.WatchOptions); err != nil {
		log.Fatalf("sauron-cli: Error retrieving latest update: %v", err)
	}
}
