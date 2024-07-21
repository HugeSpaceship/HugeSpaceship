package main

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/logger"
)

func main() {
	cfg, err := config.LoadConfig(false)
	if err != nil {
		panic(err)
	}
	logger.LoggingInit("dbmigrator", cfg)

	db.Open(cfg)
}
