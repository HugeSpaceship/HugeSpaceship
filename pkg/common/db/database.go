package db

import (
	"HugeSpaceship/pkg/common/db/migration"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
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

func open() {
	poolConfig, err := pgxpool.ParseConfig("postgres://hugespaceship:hugespaceship@localhost:5432/hugespaceship")
	if err != nil {
		panic(err.Error())
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