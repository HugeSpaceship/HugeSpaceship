package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"io"
)

func (c *Context) UploadResource(reader io.ReadCloser, contentLength int64, hash string) error {
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

	_, err = c.pool.Exec(c.ctx, "INSERT INTO resources VALUES (originaluploader,file,hash) ($1, $2, $3)")
	return err
}
