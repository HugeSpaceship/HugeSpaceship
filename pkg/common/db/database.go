package db

import (
	"HugeSpaceship/pkg/common/db/migration"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"sync"
)

type Context struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

var connection *Context

var connectOnce = sync.Once{}

func GetConnection() *Context {
	connectOnce.Do(open)
	return connection
}

func GetContext() context.Context {
	ctx := context.Background()
	conn, err := connection.pool.Acquire(ctx)
	if err != nil {
		return nil
	}
	return context.WithValue(ctx, "conn", conn)
}

func open() {
	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://hugespaceship:hugespaceship@%s:5432/hugespaceship", viper.GetString("db_host")))
	if err != nil {
		panic(err.Error())
	}

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		panic(err.Error())
	}
	if err != nil {
		panic(err.Error())
	}
	err = migration.MigrateDB(conn)
	if err != nil {
		panic(err.Error())
	}
	connection = &Context{pool: conn, ctx: context.Background()}
}
