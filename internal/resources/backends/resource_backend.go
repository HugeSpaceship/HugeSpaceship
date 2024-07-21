package backends

import (
	"github.com/spf13/viper"
	"io"
)

type ResourceBackend interface {
	InitConnection(config map[string]interface{}, viper *viper.Viper) (BackendConnection, error)
}

type BackendConnection interface {
	CanUpload() bool
	GetResource(hash string) (io.ReadCloser, int64, error)
	HasResource(hash string) (bool, error)
	HasResources(hashes []string) ([]string, error)
	DeleteResource(hash string) error
	UploadResource(hash string, r io.Reader, length int64) error
}
