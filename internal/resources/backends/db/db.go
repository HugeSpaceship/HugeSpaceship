package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log/slog"
	"strconv"
)

type Backend struct {
	pool *pgxpool.Pool
}

type dbConfig struct {
	db       string
	port     string
	hostname string
	username string
	password string
}

func decodeConfig(config map[string]string) (out *dbConfig, err error) {
	out = &dbConfig{}
	for k, v := range config {
		switch k {
		case "db":
			out.db = v
		case "port":
			out.port = v
		case "host":
			out.hostname = v
		case "username":
			out.username = v
		case "password":
			out.password = v
		default:
			return nil, fmt.Errorf("unknown db config value %s", k)
		}
	}
	return out, nil
}

func (b Backend) InitConnection(config map[string]string, globalConfig *config.Config) (backends.BackendConnection, error) {
	dbCfg, err := decodeConfig(config)
	if err != nil {
		return nil, err
	}
	if dbCfg.hostname == "" {
		dbCfg.hostname = globalConfig.Database.Host
		dbCfg.port = strconv.Itoa(int(globalConfig.Database.Port))
		dbCfg.username = globalConfig.Database.Username
		dbCfg.password = globalConfig.Database.Password
		dbCfg.db = globalConfig.Database.Database
	}

	dbOpenStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?application_name=%s",
		dbCfg.username,
		dbCfg.password,
		dbCfg.hostname,
		dbCfg.port,
		dbCfg.db,
		"HugeSpaceship+Dev+DB+Storage+Backend", // because it's a URL it needs the spaces to be escaped with + signs
	)

	pgCfg, err := pgxpool.ParseConfig(dbOpenStr) // We don't need to use the field parser because we already have all the fields
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgCfg)
	if err != nil {
		panic(err.Error())
	}

	return &BackendConnection{pool: pool}, nil
}

type BackendConnection struct {
	pool *pgxpool.Pool
}

func (b *BackendConnection) UploadResource(hash string, r io.Reader, length int64) error {
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
	_, err = tx.Exec(context.Background(), "INSERT INTO resource_files (size,file,hash) VALUES ($1, $2, $3)", written, oid, hash)
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

func (b *BackendConnection) GetResource(hash string) (io.ReadCloser, int64, error) {
	tx, err := b.pool.Begin(context.Background())
	if err != nil {
		return nil, 0, err
	}

	row := tx.QueryRow(context.Background(), "SELECT file, size FROM resource_files WHERE hash = $1", hash)
	var oid uint32
	var size int64
	err = row.Scan(&oid, &size)
	if err != nil {
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

func (b *BackendConnection) HasResource(hash string) (bool, error) {
	conn, err := b.pool.Acquire(context.Background())
	if err != nil {
		return false, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(), "SELECT count(1) FROM resource_files WHERE hash = $1", hash)
	var count uint64
	err = row.Scan(&count)

	return count > 0, err
}

const resourceCheckQuery = `
SELECT l.hash
from UNNEST($1) as l(hash)
LEFT JOIN resource_files r on l.hash = r.hash
WHERE r.hash is null;
`

func (b *BackendConnection) HasResources(hashes []string) ([]string, error) {
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

func (b *BackendConnection) DeleteResource(hash string) error {
	tx, err := b.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), "SELECT file FROM resource_files WHERE hash = $1", hash)
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

	_, err = tx.Exec(context.Background(), `DELETE FROM resource_files WHERE hash = $1`, hash)
	_ = tx.Commit(context.Background())
	return err
}

var _ backends.ResourceBackend = &Backend{}
var _ backends.BackendConnection = &BackendConnection{}
