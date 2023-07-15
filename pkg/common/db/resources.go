package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"io"
)

func (c *Context) UploadResource(reader io.ReadCloser, contentLength int64, hash string, uploader uuid.UUID) error {
	conn, err := c.pool.Acquire(c.ctx)
	if err != nil {
		return err
	}
	tx, err := conn.Begin(c.ctx)
	if err != nil {
		return err
	}
	lobs := tx.LargeObjects()
	if err != nil {
		return err
	}
	oid, err := lobs.Create(c.ctx, 0)
	if err != nil {
		return err
	}
	lob, err := lobs.Open(c.ctx, oid, pgx.LargeObjectModeWrite)
	if err != nil {
		return err
	}
	written, err := io.Copy(lob, reader)
	if err != nil {
		return err
	}
	if contentLength != written {
		log.Debug().Int64("content-length", contentLength).Int64("bytes-written", written).Msg("unexpected content length")
	}
	err = reader.Close()
	if err != nil {
		return err
	}

	_, err = c.pool.Exec(c.ctx, "INSERT INTO resources (originaluploader,size,file,hash) VALUES ($1, $2, $3, $4)", uploader, written, oid, hash)
	if err != nil {
		return err
	}

	err = tx.Commit(c.ctx)
	return err
}

func GetResource(ctx context.Context, hash string) (io.ReadSeekCloser, pgx.Tx, int64, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	row := conn.QueryRow(ctx, "SELECT oid, size FROM resources WHERE hash = $1", hash)
	var oid uint32
	var size int64
	err := row.Scan(&oid, &size)
	if err != nil {
		return nil, nil, 0, err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			log.Fatal().Err(er2).Msg("Failed to rollback broken transaction")
		}
		return nil, nil, 0, err
	}
	lobs := tx.LargeObjects()
	lob, err := lobs.Open(ctx, oid, pgx.LargeObjectModeRead)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			log.Fatal().Err(er2).Msg("Failed to rollback broken transaction")
		}
		return nil, nil, 0, err
	}

	return lob, tx, size, nil
}

func GetLevelResources(ctx context.Context, id uuid.UUID) ([]string, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)

	count := conn.QueryRow(ctx, "SELECT count(resource_hash) FROM slot_resources")
	var resCount uint
	err := count.Scan(&resCount)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, "SELECT resource_hash FROM slot_resources WHERE slot_id = $1", id)
	if err != nil {
		return nil, err
	}
	resources := make([]string, resCount)
	i := 0
	for rows.Next() {
		var resource string
		err := rows.Scan(&resource)
		if err != nil {
			return nil, err
		}
		resources[i] = resource
		i++
	}
	return resources, nil
}
