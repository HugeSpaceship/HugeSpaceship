package backends

import (
	"io"
)

type ResourceBackend interface {
	GetResource(hash string) (io.ReadCloser, error)
	HasResource(hash string) (bool, error)
	HasResources(hashes []string) ([]string, error)
	DeleteResource(hash string) error
	UploadResource(hash string, r io.Reader, length int64) error
}
