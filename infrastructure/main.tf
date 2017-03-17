variable "account_id" {
  description = "The id of the AWS in which to deploy"
}

variable "function_name" {
  description = "The name of the lambda function. Defaults to 'sauron'"
  default     = "sauron"
}

variable "bucket_name" {
  description = "The name of S3 bucket for state storage. Defaults to 'sauron-<random id>'"
  default     = ""
}

variable "key_name" {
  description = "The key to use for state storage in S3. Defaults to 'state.json'"
  default     = "state.json"
}

variable "schedule" {
  description = "A valid schedule expression to trigger the lambda function. See http://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html"
  default     = "rate(1 hour)"
}

variable "watches" {
  type    = "list"
  default = []
}

provider "aws" {
  allowed_account_ids = ["${var.account_id}"]
}

// Permissions
resource "aws_iam_role" "sauron" {
  name = "sauron-lambda-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "sauron" {
  name = "sauron-lambda-policy"
  role = "${aws_iam_role.sauron.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:*"
      ],
      "Effect": "Allow",
      "Resource": "*",
      "Sid": "AllowCloudwatchLogs"
    },
    {
      "Action": [
        "s3:ListBucket",
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Effect": "Allow",
      "Resource": [
        "${aws_s3_bucket.state.arn}",
        "${aws_s3_bucket.state.arn}/${var.key_name}"
      ],
      "Sid": "AllowS3StateReadWrite"
    },
    {
      "Action": [
        "sns:Publish"
      ],
      "Effect": "Allow",
      "Resource": [
        "${aws_sns_topic.notifications.arn}"
      ],
      "Sid": "AllowSNSPublish"
    }
  ]
}
EOF
}

// State
resource "aws_s3_bucket" "state" {
  bucket = "${coalesce(var.bucket_name, format("sauron-%10s", sha1(uuid())))}"
  acl    = "private"

  force_destroy = true

  lifecycle {
    ignore_changes = ["bucket"]
  }
}

// Notifications
resource "aws_sns_topic" "notifications" {
  name = "sauron-notifications"
}

// Lambda
resource "aws_lambda_function" "sauron" {
  function_name = "${var.function_name}"

  filename         = "${path.module}/sauron.zip"
  source_code_hash = "${base64sha256(file("${path.module}/sauron.zip"))}"

  role = "${aws_iam_role.sauron.arn}"

  runtime = "nodejs4.3"
  handler = "index.handle"

  environment {
    variables = {
      S3_BUCKET = "${aws_s3_bucket.state.id}"
      S3_KEY    = "${var.key_name}"

      SNS_TOPIC_ARN = "${aws_sns_topic.notifications.arn}"
    }
  }
}

resource "aws_lambda_permission" "allow-cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.sauron.function_name}"
  principal     = "events.amazonaws.com"
  source_arn    = "${aws_cloudwatch_event_rule.sauron.arn}"
}

resource "aws_cloudwatch_event_rule" "sauron" {
  name                = "sauron-trigger"
  description         = "Invoke sauron every ${var.schedule}"
  schedule_expression = "${var.schedule}"
}

resource "aws_cloudwatch_event_target" "sauron" {
  count = "${length(var.watches)}"

  rule = "${aws_cloudwatch_event_rule.sauron.name}"
  arn  = "${aws_lambda_function.sauron.arn}"

  input = "${jsonencode(var.watches[count.index])}"
}

// Output
output "function_name" {
  value = "${var.function_name}"
}

output "sns_topic" {
  value = "${aws_sns_topic.notifications.arn}"
}
