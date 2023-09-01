package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"io"
)

func UploadResource(ctx context.Context, reader io.ReadCloser, contentLength int64, hash string, uploader uuid.UUID) error {
	conn := ctx.Value("conn").(*pgxpool.Conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	lobs := tx.LargeObjects()
	if err != nil {
		return err
	}
	oid, err := lobs.Create(ctx, 0)
	if err != nil {
		return err
	}
	lob, err := lobs.Open(ctx, oid, pgx.LargeObjectModeWrite)
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

	_, err = conn.Exec(ctx, "INSERT INTO resources (originaluploader,size,file,hash) VALUES ($1, $2, $3, $4)", uploader, written, oid, hash)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func GetResource(ctx context.Context, hash string) (io.ReadSeekCloser, pgx.Tx, int64, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	row := conn.QueryRow(ctx, "SELECT file, size FROM resources WHERE hash = $1", hash)
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

func ResourceExists(ctx context.Context, hash string) (exists bool, err error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)

	var count int
	row := conn.QueryRow(ctx, "SELECT count(1) FROM resources WHERE hash = $1", hash)
	err = row.Scan(&count)

	if count != 0 {
		return true, err
	}
	return false, err
}
