# sauron
[![CircleCI](https://circleci.com/gh/Bowbaq/sauron.svg?style=shield&circle-token=bee68e9ea89b65e9164ef72128ecd6e8e70146aa)](https://circleci.com/gh/Bowbaq/sauron) [![GoDoc](https://godoc.org/github.com/Bowbaq/sauron?status.svg)](https://godoc.org/github.com/Bowbaq/sauron)

Watch for changes to a file in a public GitHub repository, get notifications.

## Design

Sauron is designed to be run on a schedule, whether that's a CRON job, a CloudWatch shedule or something else. Each invocation provides details about which repository to watch (and optionally which branch / path within that repository). Sauron fetches the current state and compares it against the stored state. If there is a difference, a notification is published (except on the first run).

Sauron is designed to allow for easy extension of the storage notification backends. Details below.

### Storage

Sauron currently supports the following methods for storing state:
- S3 bucket
- PostgresSQL database
- File

### Notifications
Sauron currently supports the following methods for notifying about changes:
- SNS topic
- File

## Deploy

You have a few options when it comes to deploy Sauron. Currently supported are a standlone CLI with a CRON job
and an AWS deployment using Lambda, S3 for storage and SNS for notifications.

### Standalone

Sauron comes with a standalone CLI. If you have a go toolchain, you can simply run:
```
go install github.com/Bowbaq/sauron/cmd/sauron
```

Otherwise the [releases](https://github.com/Bowbaq/sauron/releases) page will have the latest binaries for your platform.

#### Usage
```
Usage:
  sauron [OPTIONS]

notifier.sns:
      --notifier.sns.topic-arn= ARN of the SNS topic [$SNS_TOPIC_ARN]

notifier.file:
      --notifier.file.path=     path of the state file [$NOTIFY_FILE_PATH]

store.s3:
      --store.s3.bucket=        name of the bucket [$S3_BUCKET]
      --store.s3.key=           path to the key (default: state.json) [$S3_KEY]

store.postgres:
      --store.pg.datasource=    postgresql datasource (see database/sql) [$PG_DATASOURCE]

store.file:
      --store.file.path=        path of the state file (default: .sauron) [$STORE_FILE_PATH]

github:
      --github.owner=           owner of the repository [$GITHUB_OWNER]
      --github.repository=      name of the repository [$GITHUB_REPOSITORY]
      --github.branch=          branch to watch in the repository [$GITHUB_BRANCH]
      --github.path=            path to watch in the repository [$GITHUB_PATH]

Help Options:
  -h, --help                    Show this help message
```

Example CRON job:
```
0 * * * * sauron --github.owner Bowbaq --github.repository sauron --notifier.file.path sauron.changelog
```

### AWS

#### Dependencies

The AWS deployment script needs the following tools to be installed:

- [apex](https://github.com/apex/apex) - package the Go lambda into a format AWS can run
- [terraform](https://github.com/hashicorp/terraform) - manage the infrastructure
- [aws-cli](https://aws.amazon.com/cli/) - apply infrastructure tweaks that terraform cannot support

#### Configuration

1. Clone this repository in a convenient location.
    ```shell
    git clone https://github.com/Bowbaq/sauron.git
    ```
1. Create the configuration file from the template, then customize
    ```
    cp infrastructure/terraform.tfvars.template infrastructure/terraform.tfvars
    ```
    - Update the `account_id` (you can find it [here](https://console.aws.amazon.com/support/home))
    - Add your own watches to the list (by default `sauron` watches this file)

#### Deployment

1. Run the deployment script
    ```shell
    AWS_PROFILE=<profile> AWS_DEFAULT_REGION=<region> ./scripts/deploy-aws
    ```
    This will:
    - Package the lambda function
    - Create the needed infrastructure.
      > Note: This may not be free, but it should be cheap (and I'm not responsible for any cost you might incur)

Congratulations, Sauron is now publishing change events to an SNS topic. You can get the ARN by running:
```shell
(cd infrastructure; terraform output sns_topic)
```

You can get email notification by subscribing to the topic:
```shell
aws sns subscribe                                                \
  --topic-arn "$(cd infrastructure; terraform output sns_topic)" \
  --protocol email                                               \
  --notification-endpoint john.doe@example.com
```
