package db

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const insertSQL = `INSERT INTO slots (
                   id, name, description, 
                   icon, root_level, locationX, locationY, 
                   initially_locked, sub_level, lbp1only, 
                   shareable, background, level_type, 
                   min_players, max_players, move_required)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`

func InsertSlot(ctx context.Context, slot lbp_xml.SlotData) (*uuid.UUID, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			return nil, er2
		}
		return nil, err
	}

	id := uuid.New()

	_, err = tx.Exec(ctx, insertSQL, id, slot.Name, slot.Description, slot.Icon, slot.RootLevel,
		slot.Location.X, slot.Location.Y, slot.InitiallyLocked, slot.IsSubLevel, slot.IsLBP1Only, slot.Shareable,
		slot.Background, slot.LevelType, slot.MinPlayers, slot.MaxPlayers, slot.MoveRequired,
	)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			return nil, er2
		}
		return nil, err
	}

	for _, res := range slot.Resource {
		_, err := tx.Exec(ctx, "INSERT INTO slot_resources VALUES($1, $2)", id, res)
		if err != nil {
			er2 := tx.Rollback(ctx)
			if er2 != nil {
				return nil, er2
			}
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &id, err
}
