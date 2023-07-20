package model

import (
	"github.com/google/uuid"
	"time"
)

type Slot struct {
	ID              int       `db:"id"`
	Uploader        uuid.UUID `db:"uploader"`
	Name            string    `db:"name"`
	Description     string    `db:"description"`
	Icon            string    `db:"icon"`
	RootLevel       string    `db:"root_level"`
	Location        Location  `db:"location"`
	InitiallyLocked bool      `db:"initially_locked"`
	IsSubLevel      bool      `db:"sub_level"`
	IsLBP1Only      bool      `db:"lbp1_only"`
	Shareable       int       `db:"shareable"`
	Background      string    `db:"background"`
	LevelType       string    `db:"level_type"`
	MinPlayers      int       `db:"min_players"`
	MaxPlayers      int       `db:"max_players"`
	MoveRequired    bool      `db:"move_required"`
	FirstPublished  time.Time `db:"first_published"`
	LastUpdated     time.Time `db:"last_updated"`
	Domain          int       `db:"domain"`
}
