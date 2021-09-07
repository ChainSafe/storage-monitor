package s3api

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	bucketName = "storage-test"
	folder     = "test"
	sizeBytes  = 5242880
)

type s3Checker struct {
	client         *minio.Client
	serviceAddress string
}

var (
	fileKeys = []string{"test/file1", "test/file2"}
)

// New creates new instance of the S3 Checker
func New(serviceAddress string, secure bool) (*s3Checker, error) {
	cli, err := newClient(serviceAddress, secure)
	if err != nil {
		return nil, fmt.Errorf("error creating new instance of S3 API Chcker: %w", err)
	}

	return &s3Checker{
		client:         cli,
		serviceAddress: serviceAddress,
	}, nil
}

func (s *s3Checker) Info() string {
	return fmt.Sprintf("S3 API Checker for %s", s.serviceAddress)
}

func (s *s3Checker) Exec(ctx context.Context) error {
	err := s.checkTestBucket(ctx)
	if err != nil {
		return fmt.Errorf("error checking test bucket: %w", err)
	}

	err = s.testCycle(ctx)
	if err != nil {
		return fmt.Errorf("error running test cycle: %w", err)
	}
	return nil
}

// NewClient creates new minio.Client that will be used throughout all other test methods
func newClient(serviceAddress string, secure bool) (*minio.Client, error) {
	// initiate a client
	return minio.New(serviceAddress, &minio.Options{
		Creds:  credentials.NewEnvAWS(),
		Secure: secure,
		Region: "us-east-1",
	})
}

func (s *s3Checker) generateNewFiles(pwd string) error {
	err := os.Mkdir(path.Join(pwd, folder), 0750)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating test folder: %w", err)
	}

	for _, key := range fileKeys {
		filePath := path.Join(pwd, key)
		var f *os.File
		f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return fmt.Errorf("error creating file %s: %w", filePath, err)
		}
		buf := make([]byte, sizeBytes)
		_, err = rand.Read(buf)
		if err != nil {
			return fmt.Errorf("error reading from /dev/urandom for file %s: %w", filePath, err)
		}

		_, err = f.Write(buf)
		if err != nil {
			return fmt.Errorf("error writing to file %s: %w", filePath, err)
		}
		err = f.Close()
		if err != nil {
			return fmt.Errorf("error closing file %s descriptor: %w", filePath, err)
		}
	}

	return nil
}

// checkTestBucket checks if test bucket exist and if not creates one
func (s *s3Checker) checkTestBucket(ctx context.Context) error {
	// list all buckets
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking bucket existense: %w", err)
	}

	if !exists {
		// create bucket
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
	}
	return nil
}

// testCycle removes files from previous upload and uploads new ones
func (s *s3Checker) testCycle(ctx context.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %w", err)
	}

	objListOptions := minio.ListObjectsOptions{Prefix: "/", Recursive: true}
	for object := range s.client.ListObjects(ctx, bucketName, objListOptions) {
		for _, key := range fileKeys {
			bucketPath := fmt.Sprintf("/%s", key)
			if strings.EqualFold(object.Key, bucketPath) {
				obj, err := s.client.GetObject(ctx, bucketName, bucketPath, minio.GetObjectOptions{})
				if err != nil {
					return fmt.Errorf("error getting object [%s]: %w", key, err)
				}

				_, err = ioutil.ReadAll(obj)
				if err != nil {
					return fmt.Errorf("error downloading object [%s]: %w", key, err)
				}

				err = s.client.RemoveObject(ctx, bucketName, bucketPath, minio.RemoveObjectOptions{})
				if err != nil {
					return fmt.Errorf("error removing already existing object [%s]: %w", key, err)
				}
			}
		}
	}

	err = s.generateNewFiles(pwd)
	if err != nil {
		return nil
	}
	for _, key := range fileKeys {
		localPath := filepath.Join(pwd, key)
		_, err = s.client.FPutObject(
			ctx, bucketName, key, localPath, minio.PutObjectOptions{
				DisableMultipart: true,
			})
		if err != nil {
			return fmt.Errorf("error uploading object [%s]: %w", key, err)
		}
	}
	return nil
}
