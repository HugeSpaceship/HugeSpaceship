package migration

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed migrations/000_Initial_Migration.sql
var testMigration string

func TestGetMigration(t *testing.T) {
	migration, name, exists, err := GetMigration("000_Initial_Migration.sql")
	if !exists {
		t.Error("Migration not found")
	}
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, name, "000_Initial_Migration")

	expectedMigration := fmt.Sprintf(migrationTemplate, "000_Initial_Migration", testMigration, "000_Initial_Migration")

	if migration != expectedMigration {
		fmt.Println("Got: ")
		fmt.Println(migration)
		fmt.Println("Expected: ")
		fmt.Println(expectedMigration)
		t.Error("Migration does not match expected result")
	}
}
