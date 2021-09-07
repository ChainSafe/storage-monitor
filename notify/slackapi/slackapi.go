package slackapi

import (
	"context"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

type slackNotifier struct {
	client *slack.Client
}

// New creates new instance of Slack Notifier
func New() (*slackNotifier, error) {
	cli, err := newClient()
	if err != nil {
		return nil, fmt.Errorf("error creating instance of Slack Notifier: %w", err)
	}

	return &slackNotifier{
		client: cli,
	}, nil
}

// newClient creates new instance of a Slack client
func newClient() (*slack.Client, error) {
	key := os.Getenv("SLACK_API_KEY")
	if len(key) == 0 {
		return nil, fmt.Errorf("error: no slack API key provided")
	}

	debugVal := os.Getenv("DEBUG")
	debug := len(debugVal) != 0

	cli := slack.New(key, slack.OptionDebug(debug))
	return cli, nil
}

// Name is the human-readable name for the notification service
func (sn *slackNotifier) Name() string {
	return "Slack Notifier"
}

// Notify sends notification to the channel
func (sn *slackNotifier) Notify(ctx context.Context, label, msg string) error {
	chanID := os.Getenv("SLACK_CHANNEL_ID")
	if len(chanID) == 0 {
		return fmt.Errorf("error reading Slack channel ID")
	}

	attachment := slack.Attachment{
		Pretext: fmt.Sprintf("🚨 🚨 🚨 Problem with %s", label),
		Title:   msg,
	}

	_, _, _, err := sn.client.SendMessageContext(
		ctx, chanID, slack.MsgOptionUser("monitoring-bot"), slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}
