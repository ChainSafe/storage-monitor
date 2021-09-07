package monitoring

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chainsafe/storage-monitor/check"
	"github.com/chainsafe/storage-monitor/notify"
)

// Service is a main Monitoring Service instance
type Service struct {
	checkers  []check.Checker
	notifiers []notify.Notifier
}

// New creates new instance of a general monitoring Service
func New() (*Service, error) {
	checkers, err := check.NewCheckers(check.S3Checker)
	if err != nil {
		return nil, fmt.Errorf("error during checkers initialzation: %w", err)
	}

	notifiers, err := notify.NewNotifiers(notify.SlackNotifier)
	if err != nil {
		return nil, fmt.Errorf("error notifiers notifiers initialzation: %w", err)
	}

	return &Service{
		checkers:  checkers,
		notifiers: notifiers,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	timeoutVal := os.Getenv("REPEAT_EACH_MIN")
	timeout, err := strconv.Atoi(timeoutVal)
	if err != nil {
		timeout = 15
	}

	for {
		for chIdx := range s.checkers {
			go func(ctx context.Context, checker check.Checker) {
				err = checker.Exec(ctx)
				if err != nil {
					for nIdx := range s.notifiers {
						err = s.notifiers[nIdx].Notify(ctx, checker.Info(), err.Error())
						if err != nil {
							log.Printf(
								"error making notification with %s: %s", s.notifiers[nIdx].Name(), err.Error())
						}
					}
				}
			}(ctx, s.checkers[chIdx])
		}
		time.Sleep(time.Duration(timeout) * time.Minute)
	}
}
