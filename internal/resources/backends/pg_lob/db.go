package pg_lob

import (
	"context"
	"errors"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log/slog"
)

func NewBackend(pool *pgxpool.Pool) (backends.ResourceBackend, error) {
	return &PgLOBBackend{pool}, nil
}

type PgLOBBackend struct {
	pool *pgxpool.Pool
}

func (b *PgLOBBackend) UploadResource(hash string, r io.Reader, length int64) error {
	// Get a transaction
	tx, err := b.pool.Begin(context.Background())
	if err != nil {
		return err
	}

	// Create large object
	lobs := tx.LargeObjects()
	oid, err := lobs.Create(context.Background(), 0)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	// Open large object for writing
	lob, err := lobs.Open(context.Background(), oid, pgx.LargeObjectModeWrite)
	if err != nil {
		tx.Rollback(context.Background())
	}

	// Copy data into object
	written, err := io.Copy(lob, r)
	if err != nil || written != length { // Handle failure conditions by removing lob
		lob.Close()
		er2 := lobs.Unlink(context.Background(), oid)
		if er2 != nil {
			slog.Error("failed to unlink, incomplete object may be in DB", slog.Any("err", err))
		}
		tx.Rollback(context.Background())
		if written != length {
			slog.Error("object smaller than reported size, removed from DB", slog.Int64("expected", length), slog.Int64("actual", written))
			return errors.New("object smaller than reported size")
		}
		if err != nil {
			return err
		}
	}

	// Save object id to resource_files table
	lob.Close()
	_, err = tx.Exec(context.Background(), "INSERT INTO files (file,hash) VALUES ($1, $2)", oid, hash)
	if err != nil {
		er2 := lobs.Unlink(context.Background(), oid)
		if er2 != nil {
			slog.Error("failed to unlink, incomplete object may be in DB", slog.Any("err", err))
		}
		tx.Rollback(context.Background())
		return err
	}

	_ = tx.Commit(context.Background())
	return nil
}

func (b *PgLOBBackend) GetResource(hash string) (io.ReadCloser, int64, error) {
	tx, err := b.pool.Begin(context.Background())
	if err != nil {
		return nil, 0, err
	}

	row := tx.QueryRow(context.Background(), "SELECT file FROM files WHERE hash = $1", hash)
	var oid uint32
	var size int64
	err = row.Scan(&oid, &size)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, backends.ResourceNotFound
		}
		return nil, 0, err
	}

	lobs := tx.LargeObjects()
	lob, err := lobs.Open(context.Background(), oid, pgx.LargeObjectModeRead)
	if err != nil {
		tx.Rollback(context.Background())
		return nil, 0, err
	}

	return &LobCloser{
		tx, lob,
	}, size, nil
}

func (b *PgLOBBackend) HasResource(hash string) (bool, error) {
	conn, err := b.pool.Acquire(context.Background())
	if err != nil {
		return false, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(), "SELECT count(1) FROM files WHERE hash = $1", hash)
	var count uint64
	err = row.Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, backends.ResourceNotFound
		}
	}

	return count > 0, err
}

const resourceCheckQuery = `
SELECT l.hash
from UNNEST($1) as l(hash)
LEFT JOIN files r on l.hash = r.hash
WHERE r.hash is null;
`

func (b *PgLOBBackend) HasResources(hashes []string) ([]string, error) {
	conn, err := b.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	out := make([]string, 0, len(hashes))

	rows, err := conn.Query(context.Background(), resourceCheckQuery, hashes)
	if err != nil {
		return nil, err
	}
	out, err = pgx.CollectRows(rows, pgx.RowTo[string])

	return out, err
}

func (b *PgLOBBackend) DeleteResource(hash string) error {
	tx, err := b.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), "SELECT file FROM files WHERE hash = $1", hash)
	var oid uint32
	err = row.Scan(&oid)
	if err != nil {
		return err
	}

	lobs := tx.LargeObjects()
	err = lobs.Unlink(context.Background(), oid)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM files WHERE hash = $1`, hash)
	_ = tx.Commit(context.Background())
	return err
}

var _ backends.ResourceBackend = &PgLOBBackend{}
