package resources

import (
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends/pg_lob"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
)

const (
	PGLob string = "pg_lob"
)

type ResourceManager struct {
	pool    *pgxpool.Pool
	backend backends.ResourceBackend
}

func NewResourceManager(pool *pgxpool.Pool, cfg *config.Config) (*ResourceManager, error) {
	var backend backends.ResourceBackend

	switch cfg.ResourceServer.Backend {
	case PGLob:
		lobBackend, err := pg_lob.NewBackend(pool)
		if err != nil {
			return nil, err
		}
		backend = lobBackend
	default:
		return nil, fmt.Errorf("unknown backend %s", cfg.ResourceServer.Backend)
	}

	return &ResourceManager{
		pool:    pool,
		backend: backend,
	}, nil
}

func (rm *ResourceManager) GetResource(hash string) (io.ReadCloser, int64, error) {
	return rm.backend.GetResource(hash)
}

func (rm *ResourceManager) UploadResource(hash string, r io.Reader, size int64, userID uuid.UUID) error {

	return rm.backend.UploadResource(hash, r, size)
}

func (rm *ResourceManager) DeleteResource(hash string) error {
	return rm.backend.DeleteResource(hash)
}

func (rm *ResourceManager) HasResource(hash string) (bool, error) {
	return rm.backend.HasResource(hash)
}

func (rm *ResourceManager) HasResources(hashes []string) ([]string, error) {
	return rm.backend.HasResources(hashes)
}
