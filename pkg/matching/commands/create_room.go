package commands

import (
	"HugeSpaceship/pkg/matching/types"
	"HugeSpaceship/pkg/matching/types/commands"
	"context"
)

const createRoomSQL = `
INSERT INTO rooms (players, game_version, platform, room_slot) VALUES (
   platform = $1, game = $2, game_version = $3, players = $4 
);
`

func CreateRoom(ctx context.Context, room commands.CreateRoom) error {
	dbRoom := types.Room{}
}
