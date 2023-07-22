package db

import (
	"HugeSpaceship/pkg/common/model/db"
	"HugeSpaceship/pkg/common/model/lbp_xml/slot"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const insertSQL = `INSERT INTO slots (
                   name, description, 
                   icon, root_level, locationX, locationY, 
                   initially_locked, sub_level, lbp1only, 
                   shareable, background, level_type, 
                   min_players, max_players, move_required, domain, uploader, first_published, last_updated)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id;`

func InsertSlot(ctx context.Context, slot *slot.Upload, uploader uuid.UUID, domain int) (int64, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			return 0, er2
		}
		return 0, err
	}

	var id int64
	err = tx.QueryRow(ctx, insertSQL, slot.Name, slot.Description, slot.Icon, slot.RootLevel,
		slot.InitiallyLocked, slot.IsSubLevel, slot.IsLBP1Only, slot.Shareable,
		slot.Background, slot.LevelType, slot.MinPlayers, slot.MaxPlayers, slot.MoveRequired, domain, uploader,
	).Scan(&id)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			return 0, er2
		}
		return 0, err
	}

	for _, res := range slot.Resources {
		_, err := tx.Exec(ctx, "INSERT INTO slot_resources VALUES($1, $2)", id, res)
		if err != nil {
			log.Debug().Err(err).Str("hash", res).Msg("failed to insert slot resource")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}
	return id, err
}

func GetSlot(ctx context.Context, id int64) (slot.Slot, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)

	var dbSlot db.Slot

	err := pgxscan.Get(ctx, conn, &dbSlot, "SELECT * FROM slots WHERE slots.id = $1 LIMIT 1;", id)
	if err != nil {
		return slot.Slot{}, err
	}

	return slot.Slot{}, nil
}

func GetSlotsBy(ctx context.Context, by uuid.UUID, offset uint64, limit uint64) (slot.Slots[slot.SearchSlot], error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	var dbSlots []db.Slot
	err := pgxscan.Select(ctx, conn, &dbSlots, "SELECT * FROM slots WHERE uploader = $1 OFFSET $2 LIMIT $3", by, offset, limit)

	if err != nil {
		return slot.Slots[slot.SearchSlot]{}, err
	}

	slots := slot.Slots[slot.SearchSlot]{}

	slots.Total, err = GetTotalSlots(ctx)
	slots.HintStart = slots.Total - int(offset)
	return slots, nil
}

func GetTotalSlots(ctx context.Context) (int, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	row := conn.QueryRow(ctx, "SELECT COUNT(1) FROM slots;")
	var total int
	return total, row.Scan(&total)
}
