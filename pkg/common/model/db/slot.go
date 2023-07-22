package db

import (
	"HugeSpaceship/pkg/common/model/common"
	"HugeSpaceship/pkg/common/model/lbp_xml/slot"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID              int             `db:"id"`
	Uploader        uuid.UUID       `db:"uploader"`
	Name            string          `db:"name"`
	Description     string          `db:"description"`
	Icon            string          `db:"icon"`
	RootLevel       string          `db:"root_level"`
	Location        common.Location `db:"location"`
	InitiallyLocked bool            `db:"initially_locked"`
	IsSubLevel      bool            `db:"sub_level"`
	IsLBP1Only      bool            `db:"lbp1_only"`
	Shareable       int             `db:"shareable"`
	Background      string          `db:"background"`
	LevelType       string          `db:"level_type"`
	MinPlayers      int             `db:"min_players"`
	MaxPlayers      int             `db:"max_players"`
	MoveRequired    bool            `db:"move_required"`
	FirstPublished  time.Time       `db:"first_published"`
	LastUpdated     time.Time       `db:"last_updated"`
	Domain          int             `db:"domain"`
}

func SlotFromUpload(slot slot.Upload) Slot {
	return Slot{
		Name:            slot.Name,
		Description:     slot.Description,
		Icon:            slot.Icon,
		RootLevel:       slot.RootLevel,
		Location:        slot.Location,
		InitiallyLocked: slot.InitiallyLocked,
		IsSubLevel:      slot.IsSubLevel,
		IsLBP1Only:      slot.IsLBP1Only,
		Shareable:       slot.Shareable,
		Background:      slot.Background,
		LevelType:       slot.LevelType,
		MinPlayers:      slot.MinPlayers,
		MaxPlayers:      slot.MaxPlayers,
		MoveRequired:    slot.MoveRequired,
	}
}

func SearchSlotFromDB(in Slot) slot.SearchSlot {
	return slot.SearchSlot{
		Type:                "user",
		ID:                  strconv.Itoa(in.ID),
		SearchScore:         0,
		Location:            in.Location,
		Game:                0,
		Name:                "",
		Description:         "",
		RootLevel:           "",
		Icon:                "",
		InitiallyLocked:     false,
		IsSubLevel:          false,
		IsLBP1Only:          false,
		Shareable:           0,
		MinPlayers:          0,
		MaxPlayers:          0,
		HeartCount:          0,
		Thumbsup:            0,
		Thumbsdown:          0,
		AverageRating:       0,
		PlayerCount:         0,
		MatchingPlayers:     0,
		MMPick:              false,
		IsAdventurePlanet:   false,
		Ps4Only:             false,
		PlayCount:           0,
		CompletionCount:     0,
		Lbp1PlayCount:       0,
		Lbp1CompletionCount: 0,
		Lbp1UniquePlayCount: 0,
		Lbp2PlayCount:       0,
		Lbp2CompletionCount: 0,
		UniquePlayCount:     0,
		Lbp3PlayCount:       0,
		Lbp3CompletionCount: 0,
		Lbp3UniquePlayCount: 0,
	}
}
