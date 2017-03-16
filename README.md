# sauron [![CircleCI](https://circleci.com/gh/Bowbaq/sauron.svg?style=svg&circle-token=bee68e9ea89b65e9164ef72128ecd6e8e70146aa)](https://circleci.com/gh/Bowbaq/sauron)
Watch for changes to a file in a public GitHub repository, get notifications

## Deploy

### Dependencies

The deployment script needs the following tools to be installed:

- [apex](https://github.com/apex/apex) - package the Go lambda into a format AWS can run
- [terraform](https://github.com/hashicorp/terraform) - manage the infrastructure
- [aws-cli](https://aws.amazon.com/cli/) - apply infrastructure tweaks that terraform cannot support

### Configuration

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

### Deployment

1. Run the deployment script
    ```shell
    AWS_PROFILE=<profile> AWS_DEFAULT_REGION=<region> ./scripts/deploy
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
