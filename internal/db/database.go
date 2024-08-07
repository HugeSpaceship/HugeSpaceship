package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
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

func GetDSN(v *viper.Viper) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?application_name=%s",
		v.GetString("db.username"),
		v.GetString("db.password"),
		v.GetString("db.hostname"),
		v.GetString("db.port"),
		v.GetString("db.database"),
		"HugeSpaceship+Dev", // because it's a URL it needs the spaces to be escaped with + signs TODO: make this name configurable
	)
}

// Open initializes a connection to the database based on the fields in cfg.
func Open(v *viper.Viper) *pgxpool.Pool {

	dbDSN := GetDSN(v)

	dbCfg, err := pgxpool.ParseConfig(dbDSN) // We don't need to use the field parser because we already have all the fields
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
