package types

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/google/uuid"
)

type Room struct {
	Id          int             `db:"id"`
	Players     []uuid.UUID     `db:"players"`
	GameVersion common.GameType `db:"game_version"`
	Platform    common.Platform `db:"platform"`
	RoomSlot    RoomSlot        `db:"room_slot"`
}
