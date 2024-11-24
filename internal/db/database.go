package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"log/slog"
	"net/http"
	"reflect"
)

var poolConnType = reflect.TypeFor[*pgxpool.Conn]()
var noConnectionError = errors.New("no connection")

const ConnCtxKey = "db_conn"

func GetConnection(ctx context.Context) (*pgxpool.Conn, error) {
	conn := ctx.Value(ConnCtxKey)

	if conn == nil {
		return nil, noConnectionError
	}

	connType := reflect.TypeOf(conn)

	if !connType.ConvertibleTo(poolConnType) {
		return nil, noConnectionError
	}

	return conn.(*pgxpool.Conn), nil
}

func GetRequestConnection(r *http.Request) (*pgxpool.Conn, error) {
	return GetConnection(r.Context())
}

func GetDSN(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?application_name=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
		"HugeSpaceship+Dev", // because it's a URL it needs the spaces to be escaped with + signs TODO: make this name configurable
	)
}

// Open initializes a connection to the database based on the fields in cfg.
func Open(cfg *config.Config) *pgxpool.Pool {

	dbDSN := GetDSN(cfg)

	dbCfg, err := pgxpool.ParseConfig(dbDSN) // We don't need to use the field parser because we already have all the fields
	if err != nil {
		slog.Error("Failed to parse database config, please check the server configuration", "error", err)
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
