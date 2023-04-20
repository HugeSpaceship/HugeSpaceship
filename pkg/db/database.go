package db

import (
	"HugeSpaceship/pkg/db/migration"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DBOpen(connectionString string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	migration.MigrateDB(conn)
	return conn, nil
}
