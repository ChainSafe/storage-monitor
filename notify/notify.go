package notify

import (
	"context"
	"fmt"

	"github.com/chainsafe/storage-monitor/notify/slackapi"
)

const (
	// SlackNotifier represents "slack" type of notifier
	SlackNotifier = "slack"
)

// Notifier represents the basic notifier service
type Notifier interface {
	Name() string
	Notify(ctx context.Context, label, msg string) error
}

// NewNotifiers creates a new set of notifiers from provided notifier types
func NewNotifiers(notifiers ...string) ([]Notifier, error) {
	res := make([]Notifier, 0, len(notifiers))
	for _, notif := range notifiers {
		switch notif {
		case SlackNotifier:
			n, err := slackapi.New()
			if err != nil {
				return nil, fmt.Errorf("error during notifier intialization: %w", err)
			}
			res = append(res, n)
		}
	}
	return res, nil
}
