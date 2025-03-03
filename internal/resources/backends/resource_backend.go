package backends

import (
	"context"
	"io"
)

type ResourceBackend interface {
	GetResource(ctx context.Context, hash string) (io.ReadCloser, error)
	HasResource(ctx context.Context, hash string) (bool, error)
	HasResources(ctx context.Context, hashes []string) ([]string, error)
	DeleteResource(ctx context.Context, hash string) error
	UploadResource(ctx context.Context, hash string, r io.Reader, length int64) error
}
