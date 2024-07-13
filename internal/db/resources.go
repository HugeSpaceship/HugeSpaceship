package db

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"io"
)

func UploadResource(conn *pgxpool.Conn, reader io.ReadCloser, contentLength int64, hash string, uploader uuid.UUID) error {
	defer reader.Close()

	res, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%x", sha1.Sum(res)) != hash {
		return errors.New("mismatched hash")
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	lobs := tx.LargeObjects()
	if err != nil {
		return err
	}
	oid, err := lobs.Create(context.Background(), 0)
	if err != nil {
		return err
	}
	lob, err := lobs.Open(context.Background(), oid, pgx.LargeObjectModeWrite)
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
		err := lobs.Unlink(context.Background(), oid)
		if err != nil {
			return errors.Join(err, errors.New("invalid content length"))
		}
		return errors.New("invalid content length")
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO resources (originaluploader,size,file,hash) VALUES ($1, $2, $3, $4)", uploader, written, oid, hash)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	return err
}

func ResourceExists(conn *pgxpool.Conn, hash string) (bool, error) {

	row := conn.QueryRow(context.Background(), "SELECT count(1) FROM resources WHERE hash = $1", hash)
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

func CheckResources(conn *pgxpool.Conn, hashes []string) ([]string, error) {
	out := make([]string, 0, len(hashes))

	rows, err := conn.Query(context.Background(), resourceCheckQuery, hashes)
	if err != nil {
		return nil, err
	}
	out, err = pgx.CollectRows(rows, pgx.RowTo[string])

	return out, err
}

// GetResource Gets a resource from the DB as an io reader and a transaction that must be closed once the resource is no longer needed
func GetResource(conn *pgxpool.Conn, hash string) (io.ReadSeekCloser, pgx.Tx, int64, error) {

	row := conn.QueryRow(context.Background(), "SELECT file, size FROM resources WHERE hash = $1", hash)
	var oid uint32
	var size int64
	err := row.Scan(&oid, &size)
	if err != nil {
		return nil, nil, 0, err
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		er2 := tx.Rollback(context.Background())
		if er2 != nil {
			log.Fatal().Err(er2).Msg("Failed to rollback broken transaction")
		}
		return nil, nil, 0, err
	}
	lobs := tx.LargeObjects()
	lob, err := lobs.Open(context.Background(), oid, pgx.LargeObjectModeRead)
	if err != nil {
		er2 := tx.Rollback(context.Background())
		if er2 != nil {
			log.Fatal().Err(er2).Msg("Failed to rollback broken transaction")
		}
		return nil, nil, 0, err
	}

	return lob, tx, size, nil
}

func CloseResource(resource io.ReadSeekCloser, tx pgx.Tx) {
	err := resource.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close resource")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit")
	}
}
