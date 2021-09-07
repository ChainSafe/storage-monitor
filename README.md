# ChainSafe Storage Monitor

## Overview
It's a small tool that can invoke different [Checkers](https://github.com/ChainSafe/storage-monitor/blob/dev/check/check.go#L20) and if one of them fails fire all the [Notifiers](https://github.com/ChainSafe/storage-monitor/blob/dev/notify/notify.go#L16).

One service instance intended to run against one particular target. If you need to monitor different targets, fire several Monitor instances (for example via Docker). 

## Checkers

### S3 API Checker
Tries to get files from `storage-test` bucket, dowload them, remove them, generate new ones and upload back.
For it to work properly following ENV vars are required:
```dotenv
SERVICE_URL=
SECURE=
AWS_ACCESS_KEY=
AWS_SECRET_KEY=
```

## Notifiers

### Slack notifier
Slack notifier can send message to a Slack channel. For that Slack API key and channel ID needs to be provided
```dotenv
SLACK_API_KEY=
SLACK_CHANNEL_ID=
DEBUG= #This will also print notification to stdout
```

## Running

Monitor is intended to run as a Docker container, for multiple instances of a monitoring Docker Compose can be used.
`REPEAT_EACH_MIN` is an ENV var that will hold monitor repetition value in minutes. If no values if provided `15 min` is used by default.  
