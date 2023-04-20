package migration

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func MigrateDB(connection *pgxpool.Pool) error {
	rows, err := connection.Query(context.Background(), "SELECT * FROM hs_migrations")
	if err != nil {
		return err
	}
	log.Warn().Err(rows.Err()).Msg("SQL Error")
	//connection.Exec()
	return nil
}
