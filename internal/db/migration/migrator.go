package migration

import (
	"context"
	"errors"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

const migrationTableCreateSQL = `
-- Migrations table stores what migrations we've done as to not trip over ourselves
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS migrations
(
    id        serial primary key  not null,
    migration varchar(255) unique not null,
    succeeded bool                not null
);

COMMIT;
`

// MigrateDB migrates the Database using the migrations that get embedded
func MigrateDB(connection *pgxpool.Pool) error {
	startTime := time.Now()
	log.Info().Msg("Starting DB migration")

	_, err := connection.Exec(context.Background(), migrationTableCreateSQL)
	if err != nil {
		return err
	}

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
		if errors.Is(err, pgx.ErrNoRows) {
			return GetMigrationByNumber(0)
		}
		log.Error().Err(err).Msg("failed to scan migration row")
		return "", false, err
	}

	migrationNum, err := strconv.ParseInt(strings.Split(migration.Name, "_")[0], 10, 16)

	return GetMigrationByNumber(int(migrationNum) + 1)
}
