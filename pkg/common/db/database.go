package db

import (
	"HugeSpaceship/pkg/common/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type Context struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

var connection *Context

func GetContext() context.Context {
	ctx := context.Background()
	conn, err := connection.pool.Acquire(ctx)
	if err != nil {
		return nil
	}
	return context.WithValue(ctx, "conn", conn)
}

func CloseContext(ctx context.Context) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	conn.Release()
}

// Open initializes a connection to the database based on the fields in cfg.
func Open(cfg *config.Config) *pgxpool.Pool {
	dbOpenStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?application_name=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
		"HugeSpaceship+Dev", // because it's a URL it needs the spaces to be escaped with + signs
	)

	dbCfg, err := pgxpool.ParseConfig(dbOpenStr) // We don't need to use the field parser because we already have all the fields
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse DB config, check the config file")
	}

	dbCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap()) // Allows us to use UUIDs in SQL
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbCfg)
	if err != nil {
		panic(err.Error())
	}
	return pool
}
