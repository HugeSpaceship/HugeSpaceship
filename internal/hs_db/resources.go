package hs_db

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"io"
)

func UploadResource(ctx context.Context, reader io.ReadCloser, contentLength int64, hash string, uploader uuid.UUID) error {
	defer reader.Close()
	conn := ctx.Value("conn").(*pgxpool.Conn)

	res, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%x", sha1.Sum(res)) != hash {
		return errors.New("mismatched hash")
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

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

	written, err := lob.Write(res)
	if err != nil {
		return err
	}
	if contentLength != int64(written) {
		_ = lob.Truncate(0) // if it's not uploaded correctly then break
		lob.Close()
		err := lobs.Unlink(ctx, oid)
		if err != nil {
			return errors.Join(err, errors.New("invalid content length"))
		}
		return errors.New("invalid content length")
	}

	_, err = conn.Exec(ctx, "INSERT INTO resources (originaluploader,size,file,hash) VALUES ($1, $2, $3, $4)", uploader, written, oid, hash)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func ResourceExists(ctx context.Context, hash string) (bool, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	row := conn.QueryRow(ctx, "SELECT count(1) FROM resources WHERE hash = $1", hash)
	var count uint64
	err := row.Scan(&count)

	return count > 0, err
}

const resourceCheckQuery = `
SELECT l.hash
from UNNEST($1) as l(hash)
LEFT JOIN resources r on l.hash = r.hash
WHERE r.hash is null;
`

func CheckResources(ctx context.Context, hashes []string) ([]string, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	out := make([]string, 0, len(hashes))
	err := pgxscan.Select(ctx, conn, &out, resourceCheckQuery, hashes)
	return out, err
}

// GetResource Gets a resource from the DB as an io reader and a transaction that must be closed once the resource is no longer needed
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
