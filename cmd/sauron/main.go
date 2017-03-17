package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/Bowbaq/sauron"
	"github.com/Bowbaq/sauron/model"
	"github.com/Bowbaq/sauron/store"
)

var (
	version string
)

func main() {
	app := cli.NewApp()
	app.Name = "sauron"
	app.Usage = "Watch for changes in a file in a public GitHub repository, get email notifications"
	app.Authors = []cli.Author{
		{Name: "Maxime Bury"},
	}
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repository, r",
			Usage:  "The repository under watch",
			EnvVar: "GITHUB_REPOSITORY",
		},
		cli.StringFlag{
			Name:   "branch, b",
			Usage:  "Limit the watch to a particular branch",
			EnvVar: "GITHUB_BRANCH",
		},
		cli.StringFlag{
			Name:   "path, f",
			Usage:  "Limit the watch to a particular path",
			EnvVar: "GITHUB_PATH",
		},
		cli.StringFlag{
			Name:   "database, db",
			Usage:  "Database URL to store state",
			EnvVar: "DATABASE_URL",
		},
		cli.StringFlag{
			Name:   "s3-bucket",
			Usage:  "S3 bucket name to store state",
			EnvVar: "S3_BUCKET",
		},
		cli.StringFlag{
			Name:   "s3-key",
			Usage:  "S3 bucket key to store state",
			EnvVar: "S3_KEY",
		},
	}

	app.Action = func(c *cli.Context) error {
		owner, repo, err := split(c.String("repository"))
		if err != nil {
			log.Fatalf("sauron-cli: Error parsing arguments: %v", err)
		}

		s := sauron.New()
		if dbURL := c.String("database"); dbURL != "" {
			s.SetStore(store.NewPostgres(dbURL))
		}
		bucket, key := c.String("s3-bucket"), c.String("s3-key")
		if bucket != "" && key != "" {
			s.SetStore(store.NewS3(bucket, key))
		}

		opts := model.WatchOptions{
			Repository: model.Repository{
				Owner: owner,
				Name:  repo,
			},
			Branch: c.String("branch"),
			Path:   c.String("path"),
		}
		if err := s.Watch(opts); err != nil {
			log.Fatalf("sauron-cli: Error retrieving latest update: %v", err)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("sauron-cli: Unexpected error: %v", err)
	}
}

func split(ownerRepo string) (string, string, error) {
	parts := strings.SplitN(ownerRepo, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Expected <owner>/<repository>, got %s", ownerRepo)
	}

	return parts[0], parts[1], nil
}
