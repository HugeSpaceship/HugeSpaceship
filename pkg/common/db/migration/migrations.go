package migration

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//go:embed migrations
var migrations embed.FS

func ListMigrations() []string {
	migrations, err := migrations.ReadDir("migrations")
	if err != nil {
		return nil
	}
	migrationNames := make([]string, len(migrations))
	for i, migration := range migrations {
		migrationNames[i] = migration.Name()
	}
	return migrationNames
}

func GetMigrationByNumber(number int) (string, error) {
	for _, file := range ListMigrations() {
		if strings.HasPrefix(file, fmt.Sprintf("%03d", number)) {
			return GetMigration(file)
		}
	}
	return "", errors.New("no migrations found for id " + strconv.Itoa(number))
}

func GetMigration(name string) (string, error) {
	migration, err := migrations.ReadFile("migrations/" + strings.Trim(name, ".sql") + ".sql")
	if err != nil {
		return "", err
	}

	return string(migration), nil
}
