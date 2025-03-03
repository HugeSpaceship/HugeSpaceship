package resources

import (
	"context"
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/sqlc"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends/pg_lob"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends/s3"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
)

const (
	PGLob string = "pg_lob"
	S3    string = "s3"
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
	case S3:
		s3Backend, err := s3.NewBackend(cfg)
		if err != nil {
			return nil, err
		}
		backend = s3Backend
	default:
		return nil, fmt.Errorf("unknown backend %s", cfg.ResourceServer.Backend)
	}

	return &ResourceManager{
		queries: sqlc.New(pool),
		backend: backend,
	}, nil
}

func (rm *ResourceManager) GetResource(ctx context.Context, hash string) (io.ReadCloser, error) {

	return rm.backend.GetResource(ctx, hash)
}

func (rm *ResourceManager) GetResourceInfo(ctx context.Context, hash string) (sqlc.Resource, error) {
	return rm.queries.GetResource(ctx, hash)
}

func (rm *ResourceManager) UploadResource(ctx context.Context, hash string, r io.Reader, size int64, userID uuid.UUID) error {

	r, resourceType, err := utils.GetResourceType(r)
	if err != nil {
		return err
	}

	err = rm.backend.UploadResource(ctx, hash, r, size)
	if err != nil {
		return err
	}

	return rm.queries.InsertResource(ctx, sqlc.InsertResourceParams{
		Uploader: userID,
		Size:     size,
		Type:     sqlc.ResourceType(resourceType),
		Hash:     hash,
	})
}

func (rm *ResourceManager) DeleteResource(ctx context.Context, hash string) error {
	err := rm.backend.DeleteResource(ctx, hash)

	if err != nil {
		return err
	}

	return rm.queries.DeleteResource(ctx, hash)
}

// TODO: functions to handle resource overrides

func (rm *ResourceManager) HasResource(ctx context.Context, hash string) (bool, error) {
	// TODO: Resource overrides should override this.
	return rm.backend.HasResource(ctx, hash)
}

func (rm *ResourceManager) HasResources(ctx context.Context, hashes []string) ([]string, error) {
	return rm.backend.HasResources(ctx, hashes)
}
