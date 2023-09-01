package migration

import (
	_ "embed"
	"testing"
)

//go:embed migrations/000_Initial_Migration.sql
var testMigration string

func TestGetMigration(t *testing.T) {
	migration, exists, err := GetMigration("000_Initial_Migration.sql")
	if !exists {
		t.Error("Migration not found")
	}
	if err != nil {
		t.Error(err)
	}
	if migration != testMigration {
		t.Error("Migration does not match expected result")
	}
}
