package check

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/chainsafe/storage-monitor/check/s3api"
)

const (
	// S3Checker represents S3 API type of interaction checker
	S3Checker = "s3"
	// BucketsChecker = "buckets"
)

// Checker represents an uptime check that
// will be performed on provided target URL
type Checker interface {
	Info() string
	Exec(ctx context.Context) error
}

// NewCheckers creates slice of Checker(s) that will be run against provided URL
func NewCheckers(checkers ...string) ([]Checker, error) {
	serviceAddress := os.Getenv("SERVICE_URL")
	if len(serviceAddress) == 0 {
		serviceAddress = "buckets.chainsafe.io"
	}

	secure, err := strconv.ParseBool(os.Getenv("SECURE"))
	if err != nil {
		secure = true
	}

	res := make([]Checker, 0, len(checkers))
	for i := range checkers {
		switch checkers[i] {
		case S3Checker:
			s3ch, err := s3api.New(serviceAddress, secure)
			if err != nil {
				return nil, fmt.Errorf("error checker initailization: %w", err)
			}
			res = append(res, s3ch)
		}
	}
	return res, nil
}
