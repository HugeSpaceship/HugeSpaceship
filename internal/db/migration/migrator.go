package migration

import (
	"context"
	"errors"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
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
	slog.Info("Starting database migration")

	_, err := connection.Exec(context.Background(), migrationTableCreateSQL)
	if err != nil {
		return err
	}

	count := 0
	for { // While there are more migrations
		sql, migrationName, hasNext, err := nextMigration(connection)
		if !hasNext {
			break
		}
		if err != nil {
			slog.Error("Failed to get migration", "error", err)
		}
		slog.Debug("Starting migration", "migration", migrationName)
		if _, err = connection.Exec(context.Background(), sql); err != nil { // Do the migration
			slog.Error("Failed to migrate DB", "migration", migrationName, "error", err)
			return err // If something explodes, bail
		}
		slog.Debug("Migration finished", "migration", migrationName)
		count++
	}

	slog.Info("Finished database migration", "milliseconds", time.Since(startTime).Milliseconds(), "migrationCount", count)

	return nil
}

func nextMigration(conn *pgxpool.Pool) (string, string, bool, error) {
	row := conn.QueryRow(context.Background(), "SELECT * FROM migrations ORDER BY id DESC LIMIT 1")

	migration := Migration{}
	err := row.Scan(&migration.ID, &migration.Name, &migration.Succeeded)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetMigrationByNumber(0)
		}
		return "", "", false, err
	}

	migrationNum, err := strconv.ParseInt(strings.Split(migration.Name, "_")[0], 10, 16)

	return GetMigrationByNumber(int(migrationNum) + 1)
}
