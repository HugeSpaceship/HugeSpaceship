package migration

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

// MigrateDB migrates the Database using the migrations that get embedded
func MigrateDB(connection *pgxpool.Pool) error {
	startTime := time.Now()
	log.Info().Msg("Starting DB migration")
	for { // While there are more migrations
		sql, hasNext, err := nextMigration(connection)
		if !hasNext {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to migrate DB")
		}
		if _, err = connection.Exec(context.Background(), sql); err != nil { // Do the migration
			log.Fatal().Err(err).Str("migration", sql).Msg("Failed to migrate DB")
			return err // If something explodes, bail
		}
	}

	log.Info().Int("migrationMs", int(time.Since(startTime).Milliseconds())).Msg("Migration Complete")

	return nil
}

func nextMigration(conn *pgxpool.Pool) (string, bool, error) {
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
