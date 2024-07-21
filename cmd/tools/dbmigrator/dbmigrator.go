package main

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/logger"
)

func main() {
	v := config.LoadConfig(false)

	logger.LoggingInit("dbmigrator", v)

	db.Open(v)
}
