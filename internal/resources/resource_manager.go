package resources

import (
	"context"
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/sqlc"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends/pg_lob"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
)

const (
	PGLob string = "pg_lob"
)

type ResourceManager struct {
	queries *sqlc.Queries
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
		queries: sqlc.New(pool),
		backend: backend,
	}, nil
}

func (rm *ResourceManager) GetResource(hash string) (io.ReadCloser, error) {
	return rm.backend.GetResource(hash)
}

func (rm *ResourceManager) GetResourceInfo(hash string) (sqlc.Resource, error) {
	return rm.queries.GetResource(context.Background(), hash)
}

func (rm *ResourceManager) UploadResource(hash string, r io.Reader, size int64, userID uuid.UUID) error {

	r, resourceType, err := utils.GetResourceType(r)
	if err != nil {
		return err
	}

	err = rm.backend.UploadResource(hash, r, size)
	if err != nil {
		return err
	}

	return rm.queries.InsertResource(context.Background(), sqlc.InsertResourceParams{
		Uploader: userID,
		Size:     size,
		Type:     sqlc.ResourceType(resourceType),
		Hash:     hash,
	})
}

func (rm *ResourceManager) DeleteResource(hash string) error {
	err := rm.backend.DeleteResource(hash)

	if err != nil {
		return err
	}

	return rm.queries.DeleteResource(context.Background(), hash)
}

// TODO: functions to handle resource overrides

func (rm *ResourceManager) HasResource(hash string) (bool, error) {
	// TODO: Resource overrides should override this.
	return rm.backend.HasResource(hash)
}

func (rm *ResourceManager) HasResources(hashes []string) ([]string, error) {
	return rm.backend.HasResources(hashes)
}
