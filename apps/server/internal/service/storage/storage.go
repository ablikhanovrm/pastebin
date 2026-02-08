package storage

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
)

type ObjectStorage interface {
	Upload(ctx context.Context, key string, content string) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
type Service struct {
	client *s3.Client
	bucket string
	log    zerolog.Logger
}

func NewS3Storage(client *s3.Client, bucket string, log zerolog.Logger) *Service {
	return &Service{client: client, bucket: bucket, log: log}
}

func (s *Service) Upload(ctx context.Context, key string, content string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        strings.NewReader(content),
		ContentType: aws.String("text/plain; charset=utf-8"),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Get(ctx context.Context, key string) (io.ReadCloser, *int64, error) {
	res, _ := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	//if err != nil {
	//	return "", err
	//}
	//defer res.Body.Close()

	//data, err := io.ReadAll(io.LimitReader(res.Body, 5<<20)) // 5 MB

	//if err != nil {
	//	return "", err
	//}

	return res.Body, res.ContentLength, nil
}

func (s *Service) Delete(ctx context.Context, key string) error {
	return nil
}
