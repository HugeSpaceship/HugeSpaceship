package config

import (
	"testing"
)

func TestParseEnv(t *testing.T) {
	tree := ParseEnv("HS_RESOURCE_SERVER_BACKENDS")

	val, exists := tree.Get(FormatPath("pg-lob.priority"))
	t.Logf("Exists %t, Value %s", exists, *val)
}
