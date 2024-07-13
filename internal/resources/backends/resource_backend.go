package backends

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"io"
)

type ResourceBackend interface {
	InitConnection(config map[string]string, globalConfig *config.Config) (BackendConnection, error)
}

type BackendConnection interface {
	GetResource(hash string) (io.ReadCloser, int64, error)
	HasResource(hash string) (bool, error)
	HasResources(hashes []string) ([]string, error)
	DeleteResource(hash string) error
	UploadResource(hash string, r io.Reader, length int64) error
}
