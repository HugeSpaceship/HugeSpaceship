package s3

import (
	"context"
	"errors"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
)

func NewBackend(cfg *config.Config) (backends.ResourceBackend, error) {

	var endpoint *string
	if cfg.ResourceServer.Endpoint != "" {
		endpoint = &cfg.ResourceServer.Endpoint
	}

	client := s3.New(s3.Options{
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.ResourceServer.AccessKeyID, cfg.ResourceServer.AccessKeySecret, "HugeSpaceship"),
		Region:       cfg.ResourceServer.Region,
		BaseEndpoint: endpoint,
	})
	return &Backend{
		client: client,
		bucket: cfg.ResourceServer.BucketName,
	}, nil
}

type Backend struct {
	client *s3.Client
	bucket string
}

func (b *Backend) UploadResource(ctx context.Context, hash string, r io.Reader, length int64) error {
	_, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(b.bucket),
		Key:           aws.String("r/" + hash),
		Body:          r,
		ContentLength: &length,
	})
	return err
}

func (b *Backend) GetResource(ctx context.Context, hash string) (io.ReadCloser, error) {
	obj, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String("r/" + hash),
	})
	if err != nil {
		return nil, err
	}
	return obj.Body, nil
}

func (b *Backend) HasResource(ctx context.Context, hash string) (bool, error) {
	_, err := b.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String("r/" + hash),
	})

	var notFoundErr *types.NotFound
	if errors.As(err, &notFoundErr) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (b *Backend) HasResources(ctx context.Context, hashes []string) ([]string, error) {
	var existingHashes []string
	for _, hash := range hashes {
		_, err := b.HasResource(ctx, hash)
		if err != nil {
			return nil, err
		}
		existingHashes = append(existingHashes, hash)
	}
	return existingHashes, nil
}

func (b *Backend) DeleteResource(ctx context.Context, hash string) error {
	_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String("r/" + hash),
	})
	return err
}

var _ backends.ResourceBackend = &Backend{}
