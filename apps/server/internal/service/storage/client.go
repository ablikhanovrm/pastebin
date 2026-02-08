package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(cfg config.S3Config) (*s3.Client, error) {
	s3Cfg := aws.Config{
		Region: cfg.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		),
		BaseEndpoint: aws.String(cfg.Endpoint),
	}

	client := s3.NewFromConfig(s3Cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 connection failed: %w", err)
	}

	return client, nil
}
