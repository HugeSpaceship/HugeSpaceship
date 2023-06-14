package migration

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func MigrateDB(connection *pgxpool.Pool) error {
	for { // While there are more migrations
		sql, err := nextMigration(connection)
		if err != nil {
			break
		}
		if _, err = connection.Exec(context.Background(), sql); err != nil { // Do the migration
			log.Error().Err(err).Str("migration", sql).Msg("Failed to migrate")
			return err // If something explodes, bail
		}
	}

	return nil
}

func nextMigration(conn *pgxpool.Pool) (string, error) {
	row := conn.QueryRow(context.Background(), "SELECT * FROM migrations ORDER BY id DESC LIMIT 1")

	migration := Migration{}
	err := row.Scan(&migration.ID, &migration.Name, &migration.Succeeded)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan migration row")
		return GetMigrationByNumber(0)
	}

	migrationNum, err := strconv.ParseInt(strings.Split(migration.Name, "_")[0], 10, 16)

	return GetMigrationByNumber(int(migrationNum) + 1)
}
