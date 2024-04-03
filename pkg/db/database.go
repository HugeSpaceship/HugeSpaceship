package db

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

var globalPool *pgxpool.Pool

func Acquire() (*pgxpool.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return globalPool.Acquire(ctx)
}

var poolConnType = reflect.TypeFor[*pgxpool.Conn]()
var noConnectionError = errors.New("no connection")

const ConnCtxKey = "db_conn"

func GetConnection(ctx context.Context) (*pgxpool.Conn, error) {
	conn := ctx.Value(ConnCtxKey)

	connType := reflect.TypeOf(conn)

	if !connType.ConvertibleTo(poolConnType) {
		return nil, noConnectionError
	}

	return conn.(*pgxpool.Conn), nil
}

func GetRequestConnection(r *http.Request) (*pgxpool.Conn, error) {
	return GetConnection(r.Context())
}

func GetContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := globalPool.Acquire(ctx)
	if err != nil {
		slog.Error("Failed to acquire database connection", "error", err)
		return nil
	}
	return context.WithValue(context.Background(), "conn", conn)
}

func CloseContext(ctx context.Context) {
	conn := utils.GetContextValue[*pgxpool.Conn](ctx, "conn")
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
	globalPool = pool
	return pool
}
